package domain

import (
	"errors"
	"fmt"
	"sync/atomic"
)

var ErrShipmentNotFound = errors.New("shipment not found")

var shipmentSequence uint64

type ShipmentLine struct {
	ProductSKU string
	Quantity   int
}

type Shipment struct {
	ID      string
	OrderID string
	Lines   []ShipmentLine
}

func NewShipmentFromOrder(order Order) (Shipment, error) {
	if err := order.EnsureShippable(); err != nil {
		return Shipment{}, err
	}

	id := atomic.AddUint64(&shipmentSequence, 1)
	lines := make([]ShipmentLine, 0, len(order.Lines))
	for _, line := range order.Lines {
		lines = append(lines, ShipmentLine{
			ProductSKU: line.ProductSKU,
			Quantity:   line.Quantity,
		})
	}

	return Shipment{
		ID:      fmt.Sprintf("shipment-%03d", id),
		OrderID: order.ID,
		Lines:   lines,
	}, nil
}
