package ports

import "hexagonal-architecture/internal/core/domain"

type ShipmentRepository interface {
	Save(shipment domain.Shipment) error
	FindByID(id string) (domain.Shipment, error)
}
