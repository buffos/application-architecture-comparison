package returns

import "microkernel-architecture/internal/kernel"

type Service struct {
	requests    Repository
	orders      kernel.ReturnableOrderProvider
	clock       kernel.Clock
	policy      kernel.ReturnEligibilityPolicy
	idempotency kernel.IdempotencyStore
	refunds     kernel.PaymentRefund
	restock     kernel.InventoryRestock
}

func NewService(requests Repository, orders kernel.ReturnableOrderProvider, clock kernel.Clock, policy kernel.ReturnEligibilityPolicy, idempotency kernel.IdempotencyStore, refunds kernel.PaymentRefund, restock kernel.InventoryRestock) Service {
	return Service{
		requests:    requests,
		orders:      orders,
		clock:       clock,
		policy:      policy,
		idempotency: idempotency,
		refunds:     refunds,
		restock:     restock,
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
			ProductCategory:  line.ProductCategory,
			Quantity:         line.Quantity,
			UnitPrice:        line.UnitPrice,
			ReturnWindowDays: line.ReturnWindowDays,
		})
	}

	request, err := NewReturnRequest(order.OrderID, order.CustomerID, command.Reason, order.ShippedAt, s.clock.Now(), command.RequestedBy, lines)
	if err != nil {
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

func (s Service) AcceptReturn(command kernel.AcceptReturnCommand) (kernel.AcceptReturnResult, error) {
	if command.IdempotencyKey == "" {
		return kernel.AcceptReturnResult{}, kernel.ErrIdempotencyKeyRequired
	}

	if result, ok, err := s.idempotency.Find(command.IdempotencyKey); err != nil || ok {
		return kernel.AcceptReturnResult{
			ReturnRequestID: result.ReturnRequestID,
			OrderID:         result.OrderID,
			CustomerID:      result.CustomerID,
			Status:          result.Status,
			LineCount:       result.LineCount,
		}, err
	}

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
		if err := request.Reject(command.ReviewedBy, command.ReviewNote); err != nil {
			return kernel.AcceptReturnResult{}, err
		}

		if err := s.requests.Save(request); err != nil {
			return kernel.AcceptReturnResult{}, err
		}

		result := kernel.AcceptReturnResult{
			ReturnRequestID: request.ID,
			OrderID:         request.OrderID,
			CustomerID:      request.CustomerID,
			Status:          request.Status,
			LineCount:       len(request.Lines),
		}
		if err := s.idempotency.Save(command.IdempotencyKey, kernel.IdempotencyResult{
			ReturnRequestID: result.ReturnRequestID,
			OrderID:         result.OrderID,
			CustomerID:      result.CustomerID,
			Status:          result.Status,
			LineCount:       result.LineCount,
		}); err != nil {
			return kernel.AcceptReturnResult{}, err
		}
		return result, nil
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

	if err := request.Accept(command.ReviewedBy, command.ProcessedBy, command.ReviewNote); err != nil {
		return kernel.AcceptReturnResult{}, err
	}

	if err := s.requests.Save(request); err != nil {
		return kernel.AcceptReturnResult{}, err
	}

	result := kernel.AcceptReturnResult{
		ReturnRequestID: request.ID,
		OrderID:         request.OrderID,
		CustomerID:      request.CustomerID,
		Status:          request.Status,
		LineCount:       len(request.Lines),
	}
	if err := s.idempotency.Save(command.IdempotencyKey, kernel.IdempotencyResult{
		ReturnRequestID: result.ReturnRequestID,
		OrderID:         result.OrderID,
		CustomerID:      result.CustomerID,
		Status:          result.Status,
		LineCount:       result.LineCount,
	}); err != nil {
		return kernel.AcceptReturnResult{}, err
	}
	return result, nil
}

func (s Service) RejectReturn(command kernel.RejectReturnCommand) (kernel.RejectReturnResult, error) {
	if command.IdempotencyKey == "" {
		return kernel.RejectReturnResult{}, kernel.ErrIdempotencyKeyRequired
	}

	if result, ok, err := s.idempotency.Find(command.IdempotencyKey); err != nil || ok {
		return kernel.RejectReturnResult{
			ReturnRequestID: result.ReturnRequestID,
			OrderID:         result.OrderID,
			CustomerID:      result.CustomerID,
			Status:          result.Status,
			LineCount:       result.LineCount,
		}, err
	}

	request, err := s.requests.FindByID(command.ReturnRequestID)
	if err != nil {
		return kernel.RejectReturnResult{}, err
	}

	if err := request.Reject(command.ReviewedBy, command.ReviewNote); err != nil {
		return kernel.RejectReturnResult{}, err
	}

	if err := s.requests.Save(request); err != nil {
		return kernel.RejectReturnResult{}, err
	}

	result := kernel.RejectReturnResult{
		ReturnRequestID: request.ID,
		OrderID:         request.OrderID,
		CustomerID:      request.CustomerID,
		Status:          request.Status,
		LineCount:       len(request.Lines),
	}
	if err := s.idempotency.Save(command.IdempotencyKey, kernel.IdempotencyResult{
		ReturnRequestID: result.ReturnRequestID,
		OrderID:         result.OrderID,
		CustomerID:      result.CustomerID,
		Status:          result.Status,
		LineCount:       result.LineCount,
	}); err != nil {
		return kernel.RejectReturnResult{}, err
	}
	return result, nil
}

func (s Service) GetReturnRequest(query kernel.GetReturnRequestQuery) (kernel.ReturnRequestDetails, error) {
	request, err := s.requests.FindByID(query.ReturnRequestID)
	if err != nil {
		return kernel.ReturnRequestDetails{}, err
	}

	lines := make([]kernel.ReturnRequestLineDetails, 0, len(request.Lines))
	for _, line := range request.Lines {
		lines = append(lines, kernel.ReturnRequestLineDetails{
			ProductSKU:      line.ProductSKU,
			ProductCategory: line.ProductCategory,
			Quantity:        line.Quantity,
		})
	}

	return kernel.ReturnRequestDetails{
		ReturnRequestID: request.ID,
		OrderID:         request.OrderID,
		CustomerID:      request.CustomerID,
		Status:          request.Status,
		Reason:          request.Reason,
		LineCount:       len(request.Lines),
		RequestedBy:     request.RequestedBy,
		ReviewedBy:      request.ReviewedBy,
		ProcessedBy:     request.ProcessedBy,
		ReviewNote:      request.ReviewNote,
		Lines:           lines,
	}, nil
}

func (s Service) ListReturnRequests(query kernel.ListReturnRequestsQuery) ([]kernel.ReturnRequestSummary, error) {
	requests, err := s.requests.ListByStatus(query.Status)
	if err != nil {
		return nil, err
	}

	summaries := make([]kernel.ReturnRequestSummary, 0, len(requests))
	for _, request := range requests {
		summaries = append(summaries, kernel.ReturnRequestSummary{
			ReturnRequestID: request.ID,
			OrderID:         request.OrderID,
			CustomerID:      request.CustomerID,
			Status:          request.Status,
			LineCount:       len(request.Lines),
		})
	}

	return summaries, nil
}
