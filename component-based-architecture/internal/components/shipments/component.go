package shipments

import (
	"component-based-architecture/internal/components/clock"
	"fmt"
)

// Component owns shipment creation and shipment state for this lesson.
type Component struct {
	shipments map[string]Shipment
	nextID    int
	clock     clock.Reader
}

func NewComponent(clock clock.Reader) *Component {
	return &Component{shipments: make(map[string]Shipment), clock: clock}
}

func (c *Component) Create(request ShipmentRequest) (Shipment, error) {
	c.nextID++
	shipment := Shipment{
		ID: fmt.Sprintf("shipment-%03d", c.nextID), OrderID: request.OrderID, CustomerID: request.CustomerID,
		Lines: append([]ShipmentLine(nil), request.Lines...), ShippedAt: c.clock.Now(),
	}
	c.shipments[shipment.ID] = shipment
	return shipment, nil
}

var _ Creator = (*Component)(nil)
