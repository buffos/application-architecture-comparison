package entities

import (
	"fmt"
	"sync/atomic"
)

const ShipmentStatusCreated = "Created"

var shipmentSequence uint64

type ShipmentLine struct {
	SKU         string
	ProductName string
	Quantity    int
}

type Shipment struct {
	ID      string
	OrderID string
	Status  string
	Lines   []ShipmentLine
}

func NewShipmentFromPaidOrder(order Order) (Shipment, error) {
	if order.Status != OrderStatusPaid {
		return Shipment{}, ErrQuoteCannotTransition
	}

	id := atomic.AddUint64(&shipmentSequence, 1)
	lines := make([]ShipmentLine, 0, len(order.Lines))
	for _, line := range order.Lines {
		lines = append(lines, ShipmentLine{
			SKU:         line.SKU,
			ProductName: line.ProductName,
			Quantity:    line.Quantity,
		})
	}

	return Shipment{
		ID:      fmt.Sprintf("shipment-%03d", id),
		OrderID: order.ID,
		Status:  ShipmentStatusCreated,
		Lines:   lines,
	}, nil
}
