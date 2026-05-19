package application

import (
	"hexagonal-architecture/internal/core/domain"
	"hexagonal-architecture/internal/core/ports"
)

type ListReturnRequestsUseCase struct {
	returns ports.ReturnRequestRepository
}

func NewListReturnRequestsUseCase(returns ports.ReturnRequestRepository) ListReturnRequestsUseCase {
	return ListReturnRequestsUseCase{returns: returns}
}

func (uc ListReturnRequestsUseCase) Execute(status string) ([]domain.ReturnRequest, error) {
	return uc.returns.ListByStatus(status)
}
