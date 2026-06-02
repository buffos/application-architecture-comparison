package returns

import (
	"modular-monolith/internal/modules/idempotency"
	"modular-monolith/internal/modules/inventory"
	"modular-monolith/internal/modules/orders"
	"modular-monolith/internal/modules/payments"
	"modular-monolith/internal/modules/returneligibility"
	"time"
)

type ReturnableOrderSource interface {
	GetReturnableOrder(orderID string) (orders.ReturnableOrder, error)
}

type Clock interface {
	Now() time.Time
}

type RequestReturnCommand struct {
	OrderID     string
	Reason      string
	RequestedBy string
	Lines       []RequestedReturnLine
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
	IdempotencyKey  string
	ActorID         string
	ReviewNote      string
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
	idempotency idempotency.Store
	payments    payments.Refunder
	clock       Clock
}

func NewService(returns Repository, orders ReturnableOrderSource, eligibility returneligibility.Evaluator, inventory inventory.Restocker, idempotency idempotency.Store, payments payments.Refunder, clock Clock) Service {
	return Service{
		returns:     returns,
		orders:      orders,
		eligibility: eligibility,
		inventory:   inventory,
		idempotency: idempotency,
		payments:    payments,
		clock:       clock,
	}
}

func (s Service) RequestReturn(command RequestReturnCommand) (RequestReturnResult, error) {
	order, err := s.orders.GetReturnableOrder(command.OrderID)
	if err != nil {
		return RequestReturnResult{}, err
	}

	returnableOrder := ReturnableOrder{
		OrderID:    order.OrderID,
		CustomerID: order.CustomerID,
		ShippedAt:  order.ShippedAt,
		Lines:      make([]ReturnableOrderLine, 0, len(order.Lines)),
	}
	for _, line := range order.Lines {
		returnableOrder.Lines = append(returnableOrder.Lines, ReturnableOrderLine{
			ProductSKU:       line.ProductSKU,
			ProductName:      line.ProductName,
			ProductCategory:  line.ProductCategory,
			Quantity:         line.Quantity,
			ShippedQuantity:  line.ShippedQuantity,
			UnitPrice:        line.UnitPrice,
			ReturnWindowDays: line.ReturnWindowDays,
		})
	}

	returnRequest, err := NewRequestedReturnRequest(returnableOrder, command.Lines, command.Reason, s.clock.Now(), command.RequestedBy)
	if err != nil {
		return RequestReturnResult{}, err
	}

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
	if result, ok, err := s.lookupIdempotentResult(command.IdempotencyKey); err != nil || ok {
		return result, err
	}

	returnRequest, err := s.returns.FindByID(command.ReturnRequestID)
	if err != nil {
		return ReviewReturnResult{}, err
	}

	if !s.eligibility.Allows(returneligibility.ReviewRequest{
		Reason:      returnRequest.Reason,
		ShippedAt:   returnRequest.ShippedAt,
		RequestedAt: returnRequest.RequestedAt,
		Lines: func() []returneligibility.ReviewLine {
			lines := make([]returneligibility.ReviewLine, 0, len(returnRequest.Lines))
			for _, line := range returnRequest.Lines {
				lines = append(lines, returneligibility.ReviewLine{
					ReturnWindowDays: line.ReturnWindowDays,
				})
			}
			return lines
		}(),
	}) {
		if err := returnRequest.Reject(command.ActorID, command.ReviewNote); err != nil {
			return ReviewReturnResult{}, err
		}

		if err := s.returns.Save(returnRequest); err != nil {
			return ReviewReturnResult{}, err
		}

		result := ReviewReturnResult{
			ReturnRequestID: returnRequest.ID,
			OrderID:         returnRequest.OrderID,
			CustomerID:      returnRequest.CustomerID,
			Status:          returnRequest.Status,
			LineCount:       len(returnRequest.Lines),
		}
		if err := s.saveIdempotentResult(command.IdempotencyKey, result); err != nil {
			return ReviewReturnResult{}, err
		}
		return result, nil
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

	if err := returnRequest.Refund(command.ActorID, command.ActorID, command.ReviewNote); err != nil {
		return ReviewReturnResult{}, err
	}

	if err := s.returns.Save(returnRequest); err != nil {
		return ReviewReturnResult{}, err
	}

	result := ReviewReturnResult{
		ReturnRequestID: returnRequest.ID,
		OrderID:         returnRequest.OrderID,
		CustomerID:      returnRequest.CustomerID,
		Status:          returnRequest.Status,
		LineCount:       len(returnRequest.Lines),
	}
	if err := s.saveIdempotentResult(command.IdempotencyKey, result); err != nil {
		return ReviewReturnResult{}, err
	}
	return result, nil
}

func (s Service) RejectReturn(command ReviewReturnCommand) (ReviewReturnResult, error) {
	if result, ok, err := s.lookupIdempotentResult(command.IdempotencyKey); err != nil || ok {
		return result, err
	}

	returnRequest, err := s.returns.FindByID(command.ReturnRequestID)
	if err != nil {
		return ReviewReturnResult{}, err
	}

	if err := returnRequest.Reject(command.ActorID, command.ReviewNote); err != nil {
		return ReviewReturnResult{}, err
	}

	if err := s.returns.Save(returnRequest); err != nil {
		return ReviewReturnResult{}, err
	}

	result := ReviewReturnResult{
		ReturnRequestID: returnRequest.ID,
		OrderID:         returnRequest.OrderID,
		CustomerID:      returnRequest.CustomerID,
		Status:          returnRequest.Status,
		LineCount:       len(returnRequest.Lines),
	}
	if err := s.saveIdempotentResult(command.IdempotencyKey, result); err != nil {
		return ReviewReturnResult{}, err
	}
	return result, nil
}

func (s Service) lookupIdempotentResult(key string) (ReviewReturnResult, bool, error) {
	if key == "" {
		return ReviewReturnResult{}, false, nil
	}

	result, ok, err := s.idempotency.Find(key)
	if err != nil || !ok {
		return ReviewReturnResult{}, ok, err
	}

	return ReviewReturnResult{
		ReturnRequestID: result.ReturnRequestID,
		OrderID:         result.OrderID,
		CustomerID:      result.CustomerID,
		Status:          result.Status,
		LineCount:       result.LineCount,
	}, true, nil
}

func (s Service) saveIdempotentResult(key string, result ReviewReturnResult) error {
	if key == "" {
		return nil
	}

	return s.idempotency.Save(key, idempotency.Result{
		ReturnRequestID: result.ReturnRequestID,
		OrderID:         result.OrderID,
		CustomerID:      result.CustomerID,
		Status:          result.Status,
		LineCount:       result.LineCount,
	})
}
