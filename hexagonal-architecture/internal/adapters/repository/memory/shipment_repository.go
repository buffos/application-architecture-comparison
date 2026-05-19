package memory

import (
	"sync"

	"hexagonal-architecture/internal/core/domain"
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

func (r *ShipmentRepository) FindByID(id string) (domain.Shipment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	shipment, ok := r.shipments[id]
	if !ok {
		return domain.Shipment{}, domain.ErrShipmentNotFound
	}

	return shipment, nil
}

func (r *ShipmentRepository) ListByOrderID(orderID string) ([]domain.Shipment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	shipments := make([]domain.Shipment, 0, len(r.shipments))
	for _, shipment := range r.shipments {
		if orderID == "" || shipment.OrderID == orderID {
			shipments = append(shipments, shipment)
		}
	}

	return shipments, nil
}
