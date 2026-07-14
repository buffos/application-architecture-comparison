package memory

import (
	"slices"
	"sync"

	"microkernel-architecture/internal/plugins/shipments"
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

	results := make([]shipments.Shipment, 0)
	for _, shipment := range r.shipments {
		if orderID == "" || shipment.OrderID == orderID {
			results = append(results, shipment)
		}
	}

	slices.SortFunc(results, func(a shipments.Shipment, b shipments.Shipment) int {
		if a.ID < b.ID {
			return -1
		}
		if a.ID > b.ID {
			return 1
		}
		return 0
	})

	return results, nil
}
