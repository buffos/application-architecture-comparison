package application

import (
	"hexagonal-architecture/internal/core/domain"
	"hexagonal-architecture/internal/core/ports"
)

type GetShipmentUseCase struct {
	shipments ports.ShipmentRepository
}

func NewGetShipmentUseCase(shipments ports.ShipmentRepository) GetShipmentUseCase {
	return GetShipmentUseCase{shipments: shipments}
}

func (uc GetShipmentUseCase) Execute(id string) (domain.Shipment, error) {
	return uc.shipments.FindByID(id)
}
