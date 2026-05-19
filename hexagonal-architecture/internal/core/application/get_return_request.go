package application

import (
	"hexagonal-architecture/internal/core/domain"
	"hexagonal-architecture/internal/core/ports"
)

type GetReturnRequestUseCase struct {
	returns ports.ReturnRequestRepository
}

func NewGetReturnRequestUseCase(returns ports.ReturnRequestRepository) GetReturnRequestUseCase {
	return GetReturnRequestUseCase{returns: returns}
}

func (uc GetReturnRequestUseCase) Execute(id string) (domain.ReturnRequest, error) {
	return uc.returns.FindByID(id)
}
