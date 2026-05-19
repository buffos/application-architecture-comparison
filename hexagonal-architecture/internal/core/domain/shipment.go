package domain

import (
	"errors"
	"fmt"
	"sync/atomic"
)

const ShipmentStatusPartiallyShipped = "PartiallyShipped"
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

func NewShipment(order Order, requested []ShipmentLine) (Shipment, error) {
	if order.Status != OrderStatusReadyForFulfillment && order.Status != OrderStatusPartiallyShipped {
		return Shipment{}, ErrShipmentNotAllowedUntilPaymentAccepted
	}

	id := atomic.AddUint64(&shipmentSequence, 1)
	lines := make([]ShipmentLine, 0)

	if len(requested) == 0 {
		for _, line := range order.Lines {
			if line.RemainingShippableQuantity() == 0 {
				continue
			}
			lines = append(lines, ShipmentLine{
				SKU:      line.SKU,
				Quantity: line.RemainingShippableQuantity(),
			})
		}
	} else {
		lines = append(lines, requested...)
	}

	if len(lines) == 0 {
		return Shipment{}, ErrShipmentLineInvalid
	}

	status := ShipmentStatusShipped
	for _, line := range order.Lines {
		requestedQuantity := 0
		for _, shippedLine := range lines {
			if shippedLine.SKU == line.SKU {
				requestedQuantity += shippedLine.Quantity
			}
		}
		if requestedQuantity < line.RemainingShippableQuantity() {
			status = ShipmentStatusPartiallyShipped
			break
		}
	}

	return Shipment{
		ID:      fmt.Sprintf("ship-%03d", id),
		OrderID: order.ID,
		Status:  status,
		Lines:   lines,
	}, nil
}
