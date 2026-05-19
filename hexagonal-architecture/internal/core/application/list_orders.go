package application

import (
	"hexagonal-architecture/internal/core/domain"
	"hexagonal-architecture/internal/core/ports"
)

type ListOrdersUseCase struct {
	orders ports.OrderRepository
}

func NewListOrdersUseCase(orders ports.OrderRepository) ListOrdersUseCase {
	return ListOrdersUseCase{orders: orders}
}

func (uc ListOrdersUseCase) Execute(status string) ([]domain.Order, error) {
	return uc.orders.ListByStatus(status)
}
