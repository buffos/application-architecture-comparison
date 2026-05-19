package application

import (
	"hexagonal-architecture/internal/core/domain"
	"hexagonal-architecture/internal/core/ports"
)

type ListShipmentsUseCase struct {
	shipments ports.ShipmentRepository
}

func NewListShipmentsUseCase(shipments ports.ShipmentRepository) ListShipmentsUseCase {
	return ListShipmentsUseCase{shipments: shipments}
}

func (uc ListShipmentsUseCase) Execute(orderID string) ([]domain.Shipment, error) {
	return uc.shipments.ListByOrderID(orderID)
}
