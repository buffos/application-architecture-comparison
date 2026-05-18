package application

import (
	"hexagonal-architecture/internal/core/domain"
	"hexagonal-architecture/internal/core/ports"
)

type CancelOrderUseCase struct {
	orders    ports.OrderRepository
	inventory ports.InventoryRelease
}

func NewCancelOrderUseCase(orders ports.OrderRepository, inventory ports.InventoryRelease) CancelOrderUseCase {
	return CancelOrderUseCase{
		orders:    orders,
		inventory: inventory,
	}
}

func (uc CancelOrderUseCase) Execute(orderID string) (domain.Order, error) {
	order, err := uc.orders.FindByID(orderID)
	if err != nil {
		return domain.Order{}, err
	}

	if err := order.Cancel(); err != nil {
		return domain.Order{}, err
	}

	reservations := make([]domain.ReservationLine, 0, len(order.Lines))
	for _, line := range order.Lines {
		reservations = append(reservations, domain.ReservationLine{
			SKU:      line.SKU,
			Quantity: line.Quantity,
		})
	}

	if err := uc.inventory.Release(reservations); err != nil {
		return domain.Order{}, err
	}

	if err := uc.orders.Save(order); err != nil {
		return domain.Order{}, err
	}

	return order, nil
}
