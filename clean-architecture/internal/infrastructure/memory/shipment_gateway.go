package memory

import (
	"sort"
	"sync"

	"clean-architecture/internal/entities"
)

type ShipmentGateway struct {
	mu        sync.RWMutex
	shipments map[string]entities.Shipment
}

func NewShipmentGateway() *ShipmentGateway {
	return &ShipmentGateway{
		shipments: make(map[string]entities.Shipment),
	}
}

func (g *ShipmentGateway) Save(shipment entities.Shipment) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.shipments[shipment.ID] = shipment
	return nil
}

func (g *ShipmentGateway) FindByID(id string) (entities.Shipment, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	shipment, ok := g.shipments[id]
	if !ok {
		return entities.Shipment{}, entities.ErrQuoteNotFound
	}

	return shipment, nil
}

func (g *ShipmentGateway) ListByOrderID(orderID string) ([]entities.Shipment, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	shipments := make([]entities.Shipment, 0, len(g.shipments))
	for _, shipment := range g.shipments {
		if orderID != "" && shipment.OrderID != orderID {
			continue
		}

		shipments = append(shipments, shipment)
	}

	sort.Slice(shipments, func(i int, j int) bool {
		return shipments[i].ID < shipments[j].ID
	})

	return shipments, nil
}
