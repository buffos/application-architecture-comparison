package returns

import (
	"modular-monolith/internal/modules/inventory"
	"modular-monolith/internal/modules/orders"
	"modular-monolith/internal/modules/payments"
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

type Service struct {
	returns   Repository
	orders    ReturnableOrderSource
	inventory inventory.Restocker
	payments  payments.Refunder
}

func NewService(returns Repository, orders ReturnableOrderSource, inventory inventory.Restocker, payments payments.Refunder) Service {
	return Service{
		returns:   returns,
		orders:    orders,
		inventory: inventory,
		payments:  payments,
	}
}

func (s Service) RequestReturn(command RequestReturnCommand) (RequestReturnResult, error) {
	order, err := s.orders.GetReturnableOrder(command.OrderID)
	if err != nil {
		return RequestReturnResult{}, err
	}

	totalAmount := 0
	for _, line := range order.Lines {
		totalAmount += line.Quantity * line.UnitPrice
	}

	if err := s.payments.Refund(payments.RefundRequest{
		OrderID:    order.OrderID,
		CustomerID: order.CustomerID,
		Amount:     totalAmount,
		Reason:     command.Reason,
	}); err != nil {
		return RequestReturnResult{}, err
	}

	restockItems := make([]inventory.RestockItem, 0, len(order.Lines))
	for _, line := range order.Lines {
		restockItems = append(restockItems, inventory.RestockItem{
			ProductSKU: line.ProductSKU,
			Quantity:   line.Quantity,
		})
	}

	if err := s.inventory.Restock(restockItems); err != nil {
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

	returnRequest := NewRefundedReturnRequest(ReturnableOrder{
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
