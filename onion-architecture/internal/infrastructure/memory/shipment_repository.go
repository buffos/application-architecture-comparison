package memory

import (
	"sync"

	"onion-architecture/internal/domain"
)

type ShipmentRepository struct {
	mu        sync.RWMutex
	shipments map[string]domain.Shipment
}

func NewShipmentRepository() *ShipmentRepository {
	return &ShipmentRepository{
		shipments: make(map[string]domain.Shipment),
	}
}

func (r *ShipmentRepository) Save(shipment domain.Shipment) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.shipments[shipment.ID] = shipment
	return nil
}
