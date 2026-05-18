package application

import (
	"hexagonal-architecture/internal/core/domain"
	"hexagonal-architecture/internal/core/ports"
)

type CreateShipmentUseCase struct {
	orders    ports.OrderRepository
	shipments ports.ShipmentRepository
	inventory ports.InventoryConsumption
}

func NewCreateShipmentUseCase(orders ports.OrderRepository, shipments ports.ShipmentRepository, inventory ports.InventoryConsumption) CreateShipmentUseCase {
	return CreateShipmentUseCase{
		orders:    orders,
		shipments: shipments,
		inventory: inventory,
	}
}

func (uc CreateShipmentUseCase) Execute(orderID string) (domain.Shipment, error) {
	order, err := uc.orders.FindByID(orderID)
	if err != nil {
		return domain.Shipment{}, err
	}

	shipment, err := domain.NewShipment(order)
	if err != nil {
		return domain.Shipment{}, err
	}

	reservations := make([]domain.ReservationLine, 0, len(order.Lines))
	for _, line := range order.Lines {
		reservations = append(reservations, domain.ReservationLine{
			SKU:      line.SKU,
			Quantity: line.Quantity,
		})
	}

	if err := uc.inventory.Consume(reservations); err != nil {
		return domain.Shipment{}, err
	}

	if err := order.MarkShipped(); err != nil {
		return domain.Shipment{}, err
	}

	if err := uc.orders.Save(order); err != nil {
		return domain.Shipment{}, err
	}

	if err := uc.shipments.Save(shipment); err != nil {
		return domain.Shipment{}, err
	}

	return shipment, nil
}
