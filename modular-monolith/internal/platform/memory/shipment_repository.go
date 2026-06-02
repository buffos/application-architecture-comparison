package memory

import (
	"sync"

	"modular-monolith/internal/modules/shipments"
)

type ShipmentRepository struct {
	mu        sync.RWMutex
	shipments map[string]shipments.Shipment
}

func NewShipmentRepository() *ShipmentRepository {
	return &ShipmentRepository{
		shipments: make(map[string]shipments.Shipment),
	}
}

func (r *ShipmentRepository) Save(shipment shipments.Shipment) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.shipments[shipment.ID] = shipment
	return nil
}

func (r *ShipmentRepository) FindByID(id string) (shipments.Shipment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	shipment, ok := r.shipments[id]
	if !ok {
		return shipments.Shipment{}, shipments.ErrShipmentNotFound
	}

	return shipment, nil
}

func (r *ShipmentRepository) ListByOrderID(orderID string) ([]shipments.Shipment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]shipments.Shipment, 0, len(r.shipments))
	for _, shipment := range r.shipments {
		if orderID == "" || shipment.OrderID == orderID {
			list = append(list, shipment)
		}
	}

	return list, nil
}
