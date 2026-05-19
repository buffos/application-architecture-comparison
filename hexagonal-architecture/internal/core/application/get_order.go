package application

import (
	"hexagonal-architecture/internal/core/domain"
	"hexagonal-architecture/internal/core/ports"
)

type GetOrderUseCase struct {
	orders ports.OrderRepository
}

func NewGetOrderUseCase(orders ports.OrderRepository) GetOrderUseCase {
	return GetOrderUseCase{orders: orders}
}

func (uc GetOrderUseCase) Execute(id string) (domain.Order, error) {
	return uc.orders.FindByID(id)
}
