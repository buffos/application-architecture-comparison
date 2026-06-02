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
	ProductSKU  string
	ProductName string
	Quantity    int
}

type ShipmentRequest struct {
	OrderID    string
	CustomerID string
	Lines      []ShipmentLine
}

func NewShipment(request ShipmentRequest) Shipment {
	id := atomic.AddUint64(&shipmentSequence, 1)

	lines := make([]ShipmentLine, 0, len(request.Lines))
	lines = append(lines, request.Lines...)

	return Shipment{
		ID:         fmt.Sprintf("shipment-%03d", id),
		OrderID:    request.OrderID,
		CustomerID: request.CustomerID,
		Lines:      lines,
	}
}
