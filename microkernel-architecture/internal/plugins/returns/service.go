package returns

import "microkernel-architecture/internal/kernel"

type Service struct {
	requests Repository
	orders   kernel.ReturnableOrderProvider
	refunds  kernel.PaymentRefund
	restock  kernel.InventoryRestock
}

func NewService(requests Repository, orders kernel.ReturnableOrderProvider, refunds kernel.PaymentRefund, restock kernel.InventoryRestock) Service {
	return Service{
		requests: requests,
		orders:   orders,
		refunds:  refunds,
		restock:  restock,
	}
}

func (s Service) RequestReturn(command kernel.RequestReturnCommand) (kernel.RequestReturnResult, error) {
	order, err := s.orders.GetReturnableOrder(command.OrderID)
	if err != nil {
		return kernel.RequestReturnResult{}, err
	}

	lines := make([]ReturnLine, 0, len(order.Lines))
	items := make([]kernel.InventoryReservationItem, 0, len(order.Lines))
	for _, line := range order.Lines {
		lines = append(lines, ReturnLine{
			ProductSKU: line.ProductSKU,
			Quantity:   line.Quantity,
			UnitPrice:  line.UnitPrice,
		})
		items = append(items, kernel.InventoryReservationItem{
			ProductSKU: line.ProductSKU,
			Quantity:   line.Quantity,
		})
	}

	request := NewReturnRequest(order.OrderID, order.CustomerID, command.Reason, lines)
	if err := s.refunds.Refund(order.OrderID, request.TotalAmount()); err != nil {
		return kernel.RequestReturnResult{}, err
	}

	if err := s.restock.Restock(items); err != nil {
		return kernel.RequestReturnResult{}, err
	}

	if err := s.requests.Save(request); err != nil {
		return kernel.RequestReturnResult{}, err
	}

	return kernel.RequestReturnResult{
		ReturnRequestID: request.ID,
		OrderID:         request.OrderID,
		CustomerID:      request.CustomerID,
		Status:          request.Status,
		LineCount:       len(request.Lines),
	}, nil
}
