package shipments

import (
	"errors"
	"fmt"
	"sync/atomic"
)

var ErrShipmentNotFound = errors.New("shipment not found")

var shipmentSequence uint64

type Shipment struct {
	ID         string
	OrderID    string
	CustomerID string
	Lines      []ShipmentLine
}

type ShipmentLine struct {
	ProductSKU string
	Quantity   int
}

func NewShipment(orderID string, customerID string, lines []ShipmentLine) Shipment {
	id := atomic.AddUint64(&shipmentSequence, 1)

	return Shipment{
		ID:         fmt.Sprintf("shipment-%03d", id),
		OrderID:    orderID,
		CustomerID: customerID,
		Lines:      lines,
	}
}
