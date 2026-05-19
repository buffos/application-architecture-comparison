package application

import (
	"hexagonal-architecture/internal/core/domain"
	"hexagonal-architecture/internal/core/ports"
)

type CompleteRefundUseCase struct {
	returns   ports.ReturnRequestRepository
	refunds   ports.RefundGateway
	inventory ports.InventoryRestock
}

func NewCompleteRefundUseCase(returns ports.ReturnRequestRepository, refunds ports.RefundGateway, inventory ports.InventoryRestock) CompleteRefundUseCase {
	return CompleteRefundUseCase{
		returns:   returns,
		refunds:   refunds,
		inventory: inventory,
	}
}

func (uc CompleteRefundUseCase) Execute(returnRequestID string) (domain.ReturnRequest, error) {
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

	if err := request.MarkRefunded(); err != nil {
		return domain.ReturnRequest{}, err
	}

	if err := uc.returns.Save(request); err != nil {
		return domain.ReturnRequest{}, err
	}

	return request, nil
}
