package memory

import (
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
