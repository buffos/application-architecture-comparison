package shipments

import "fmt"

// Component owns shipment creation and shipment state for this lesson.
type Component struct {
	shipments map[string]Shipment
	nextID    int
}

func NewComponent() *Component {
	return &Component{shipments: make(map[string]Shipment)}
}

func (c *Component) Create(request ShipmentRequest) (Shipment, error) {
	c.nextID++
	shipment := Shipment{
		ID: fmt.Sprintf("shipment-%03d", c.nextID), OrderID: request.OrderID, CustomerID: request.CustomerID,
		Lines: append([]ShipmentLine(nil), request.Lines...),
	}
	c.shipments[shipment.ID] = shipment
	return shipment, nil
}

var _ Creator = (*Component)(nil)
