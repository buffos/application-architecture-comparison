package returns

import "microkernel-architecture/internal/kernel"

type Service struct {
	requests Repository
	orders   kernel.ReturnableOrderProvider
	clock    kernel.Clock
	policy   kernel.ReturnEligibilityPolicy
	refunds  kernel.PaymentRefund
	restock  kernel.InventoryRestock
}

func NewService(requests Repository, orders kernel.ReturnableOrderProvider, clock kernel.Clock, policy kernel.ReturnEligibilityPolicy, refunds kernel.PaymentRefund, restock kernel.InventoryRestock) Service {
	return Service{
		requests: requests,
		orders:   orders,
		clock:    clock,
		policy:   policy,
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
	for _, line := range order.Lines {
		lines = append(lines, ReturnLine{
			ProductSKU:       line.ProductSKU,
			Quantity:         line.Quantity,
			UnitPrice:        line.UnitPrice,
			ReturnWindowDays: line.ReturnWindowDays,
		})
	}

	request := NewReturnRequest(order.OrderID, order.CustomerID, command.Reason, order.ShippedAt, s.clock.Now(), lines)
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

func (s Service) AcceptReturn(command kernel.AcceptReturnCommand) (kernel.AcceptReturnResult, error) {
	request, err := s.requests.FindByID(command.ReturnRequestID)
	if err != nil {
		return kernel.AcceptReturnResult{}, err
	}

	if !s.policy.Allows(kernel.ReturnEligibilityReview{
		Reason:      request.Reason,
		ShippedAt:   request.ShippedAt,
		RequestedAt: request.RequestedAt,
		Lines: func() []kernel.ReturnEligibilityLine {
			lines := make([]kernel.ReturnEligibilityLine, 0, len(request.Lines))
			for _, line := range request.Lines {
				lines = append(lines, kernel.ReturnEligibilityLine{
					ReturnWindowDays: line.ReturnWindowDays,
				})
			}
			return lines
		}(),
	}) {
		if err := request.Reject(); err != nil {
			return kernel.AcceptReturnResult{}, err
		}

		if err := s.requests.Save(request); err != nil {
			return kernel.AcceptReturnResult{}, err
		}

		return kernel.AcceptReturnResult{
			ReturnRequestID: request.ID,
			OrderID:         request.OrderID,
			CustomerID:      request.CustomerID,
			Status:          request.Status,
			LineCount:       len(request.Lines),
		}, nil
	}

	items := make([]kernel.InventoryReservationItem, 0, len(request.Lines))
	for _, line := range request.Lines {
		items = append(items, kernel.InventoryReservationItem{
			ProductSKU: line.ProductSKU,
			Quantity:   line.Quantity,
		})
	}

	if err := s.refunds.Refund(request.OrderID, request.TotalAmount()); err != nil {
		return kernel.AcceptReturnResult{}, err
	}

	if err := s.restock.Restock(items); err != nil {
		return kernel.AcceptReturnResult{}, err
	}

	if err := request.Accept(); err != nil {
		return kernel.AcceptReturnResult{}, err
	}

	if err := s.requests.Save(request); err != nil {
		return kernel.AcceptReturnResult{}, err
	}

	return kernel.AcceptReturnResult{
		ReturnRequestID: request.ID,
		OrderID:         request.OrderID,
		CustomerID:      request.CustomerID,
		Status:          request.Status,
		LineCount:       len(request.Lines),
	}, nil
}

func (s Service) RejectReturn(command kernel.RejectReturnCommand) (kernel.RejectReturnResult, error) {
	request, err := s.requests.FindByID(command.ReturnRequestID)
	if err != nil {
		return kernel.RejectReturnResult{}, err
	}

	if err := request.Reject(); err != nil {
		return kernel.RejectReturnResult{}, err
	}

	if err := s.requests.Save(request); err != nil {
		return kernel.RejectReturnResult{}, err
	}

	return kernel.RejectReturnResult{
		ReturnRequestID: request.ID,
		OrderID:         request.OrderID,
		CustomerID:      request.CustomerID,
		Status:          request.Status,
		LineCount:       len(request.Lines),
	}, nil
}
