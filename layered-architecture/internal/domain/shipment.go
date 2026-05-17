package domain

import (
	"errors"
	"fmt"
	"sync/atomic"
)

const ShipmentStatusShipped = "Shipped"

var shipmentSequence uint64

var ErrShipmentNotAllowedUntilPaymentAccepted = errors.New("shipment is not allowed until payment is accepted")
var ErrShipmentNotFound = errors.New("shipment not found")

type ShipmentLine struct {
	SKU      string
	Quantity int
}

type Shipment struct {
	ID      string
	OrderID string
	Status  string
	Lines   []ShipmentLine
}

func NewShipment(order Order) (Shipment, error) {
	if order.Status != OrderStatusReadyForFulfillment {
		return Shipment{}, ErrShipmentNotAllowedUntilPaymentAccepted
	}

	id := atomic.AddUint64(&shipmentSequence, 1)
	lines := make([]ShipmentLine, 0, len(order.Lines))

	for _, line := range order.Lines {
		lines = append(lines, ShipmentLine{
			SKU:      line.SKU,
			Quantity: line.Quantity,
		})
	}

	return Shipment{
		ID:      fmt.Sprintf("ship-%03d", id),
		OrderID: order.ID,
		Status:  ShipmentStatusShipped,
		Lines:   lines,
	}, nil
}
