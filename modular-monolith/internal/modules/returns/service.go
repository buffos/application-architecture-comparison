package returns

import (
	"modular-monolith/internal/modules/inventory"
	"modular-monolith/internal/modules/orders"
	"modular-monolith/internal/modules/payments"
	"modular-monolith/internal/modules/returneligibility"
)

type ReturnableOrderSource interface {
	GetReturnableOrder(orderID string) (orders.ReturnableOrder, error)
}

type RequestReturnCommand struct {
	OrderID string
	Reason  string
}

type RequestReturnResult struct {
	ReturnRequestID string
	OrderID         string
	CustomerID      string
	Status          string
	LineCount       int
}

type ReviewReturnCommand struct {
	ReturnRequestID string
}

type ReviewReturnResult struct {
	ReturnRequestID string
	OrderID         string
	CustomerID      string
	Status          string
	LineCount       int
}

type Service struct {
	returns     Repository
	orders      ReturnableOrderSource
	eligibility returneligibility.Evaluator
	inventory   inventory.Restocker
	payments    payments.Refunder
}

func NewService(returns Repository, orders ReturnableOrderSource, eligibility returneligibility.Evaluator, inventory inventory.Restocker, payments payments.Refunder) Service {
	return Service{
		returns:     returns,
		orders:      orders,
		eligibility: eligibility,
		inventory:   inventory,
		payments:    payments,
	}
}

func (s Service) RequestReturn(command RequestReturnCommand) (RequestReturnResult, error) {
	order, err := s.orders.GetReturnableOrder(command.OrderID)
	if err != nil {
		return RequestReturnResult{}, err
	}

	lines := make([]ReturnableOrderLine, 0, len(order.Lines))
	for _, line := range order.Lines {
		lines = append(lines, ReturnableOrderLine{
			ProductSKU:      line.ProductSKU,
			ProductName:     line.ProductName,
			ProductCategory: line.ProductCategory,
			Quantity:        line.Quantity,
			UnitPrice:       line.UnitPrice,
		})
	}

	returnRequest := NewRequestedReturnRequest(ReturnableOrder{
		OrderID:    order.OrderID,
		CustomerID: order.CustomerID,
		Lines:      lines,
	}, command.Reason)

	if err := s.returns.Save(returnRequest); err != nil {
		return RequestReturnResult{}, err
	}

	return RequestReturnResult{
		ReturnRequestID: returnRequest.ID,
		OrderID:         returnRequest.OrderID,
		CustomerID:      returnRequest.CustomerID,
		Status:          returnRequest.Status,
		LineCount:       len(returnRequest.Lines),
	}, nil
}

func (s Service) AcceptReturn(command ReviewReturnCommand) (ReviewReturnResult, error) {
	returnRequest, err := s.returns.FindByID(command.ReturnRequestID)
	if err != nil {
		return ReviewReturnResult{}, err
	}

	if !s.eligibility.Allows(returneligibility.ReviewRequest{
		Reason: returnRequest.Reason,
	}) {
		if err := returnRequest.Reject(); err != nil {
			return ReviewReturnResult{}, err
		}

		if err := s.returns.Save(returnRequest); err != nil {
			return ReviewReturnResult{}, err
		}

		return ReviewReturnResult{
			ReturnRequestID: returnRequest.ID,
			OrderID:         returnRequest.OrderID,
			CustomerID:      returnRequest.CustomerID,
			Status:          returnRequest.Status,
			LineCount:       len(returnRequest.Lines),
		}, nil
	}

	totalAmount := 0
	restockItems := make([]inventory.RestockItem, 0, len(returnRequest.Lines))
	for _, line := range returnRequest.Lines {
		totalAmount += line.Quantity * line.UnitPrice
		restockItems = append(restockItems, inventory.RestockItem{
			ProductSKU: line.ProductSKU,
			Quantity:   line.Quantity,
		})
	}

	if err := s.payments.Refund(payments.RefundRequest{
		OrderID:    returnRequest.OrderID,
		CustomerID: returnRequest.CustomerID,
		Amount:     totalAmount,
		Reason:     returnRequest.Reason,
	}); err != nil {
		return ReviewReturnResult{}, err
	}

	if err := s.inventory.Restock(restockItems); err != nil {
		return ReviewReturnResult{}, err
	}

	if err := returnRequest.Refund(); err != nil {
		return ReviewReturnResult{}, err
	}

	if err := s.returns.Save(returnRequest); err != nil {
		return ReviewReturnResult{}, err
	}

	return ReviewReturnResult{
		ReturnRequestID: returnRequest.ID,
		OrderID:         returnRequest.OrderID,
		CustomerID:      returnRequest.CustomerID,
		Status:          returnRequest.Status,
		LineCount:       len(returnRequest.Lines),
	}, nil
}

func (s Service) RejectReturn(command ReviewReturnCommand) (ReviewReturnResult, error) {
	returnRequest, err := s.returns.FindByID(command.ReturnRequestID)
	if err != nil {
		return ReviewReturnResult{}, err
	}

	if err := returnRequest.Reject(); err != nil {
		return ReviewReturnResult{}, err
	}

	if err := s.returns.Save(returnRequest); err != nil {
		return ReviewReturnResult{}, err
	}

	return ReviewReturnResult{
		ReturnRequestID: returnRequest.ID,
		OrderID:         returnRequest.OrderID,
		CustomerID:      returnRequest.CustomerID,
		Status:          returnRequest.Status,
		LineCount:       len(returnRequest.Lines),
	}, nil
}
