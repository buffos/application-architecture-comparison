package application

import (
	"hexagonal-architecture/internal/core/domain"
	"hexagonal-architecture/internal/core/ports"
)

type CompleteRefundUseCase struct {
	returns   ports.ReturnRequestRepository
	refunds   ports.RefundGateway
	inventory ports.InventoryRestock
	keys      ports.IdempotencyStore
}

func NewCompleteRefundUseCase(returns ports.ReturnRequestRepository, refunds ports.RefundGateway, inventory ports.InventoryRestock, keys ports.IdempotencyStore) CompleteRefundUseCase {
	return CompleteRefundUseCase{
		returns:   returns,
		refunds:   refunds,
		inventory: inventory,
		keys:      keys,
	}
}

func (uc CompleteRefundUseCase) Execute(returnRequestID, processedBy, idempotencyKey string) (domain.ReturnRequest, error) {
	seen, err := uc.keys.Seen("complete-refund", idempotencyKey)
	if err != nil {
		return domain.ReturnRequest{}, err
	}
	if seen {
		storedID, err := uc.keys.ResourceID("complete-refund", idempotencyKey)
		if err != nil {
			return domain.ReturnRequest{}, err
		}
		return uc.returns.FindByID(storedID)
	}

	request, err := uc.returns.FindByID(returnRequestID)
	if err != nil {
		return domain.ReturnRequest{}, err
	}

	ok, err := uc.refunds.Refund(request)
	if err != nil {
		return domain.ReturnRequest{}, err
	}

	if !ok {
		return domain.ReturnRequest{}, domain.ErrRefundFailed
	}

	if err := request.MarkRefunded(processedBy); err != nil {
		return domain.ReturnRequest{}, err
	}

	lines := make([]domain.ReservationLine, 0, len(request.Lines))
	for _, line := range request.Lines {
		lines = append(lines, domain.ReservationLine{
			SKU:      line.SKU,
			Quantity: line.Quantity,
		})
	}

	if err := uc.inventory.Restock(lines); err != nil {
		return domain.ReturnRequest{}, err
	}

	if err := uc.returns.Save(request); err != nil {
		return domain.ReturnRequest{}, err
	}

	if err := uc.keys.Remember("complete-refund", idempotencyKey, request.ID); err != nil {
		return domain.ReturnRequest{}, err
	}

	return request, nil
}
