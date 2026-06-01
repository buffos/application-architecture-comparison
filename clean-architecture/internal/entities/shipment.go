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

func NewShipmentFromOrder(order Order, lines []ShipmentLine) (Shipment, error) {
	if order.Status != OrderStatusPaid && order.Status != OrderStatusPartiallyShipped {
		return Shipment{}, ErrQuoteCannotTransition
	}

	id := atomic.AddUint64(&shipmentSequence, 1)
	shipmentLines := lines
	if len(shipmentLines) == 0 {
		shipmentLines = make([]ShipmentLine, 0, len(order.Lines))
		for _, line := range order.Lines {
			remaining := line.Quantity - line.ShippedQuantity
			if remaining <= 0 {
				continue
			}

			shipmentLines = append(shipmentLines, ShipmentLine{
				SKU:         line.SKU,
				ProductName: line.ProductName,
				Quantity:    remaining,
			})
		}
	}

	if len(shipmentLines) == 0 {
		return Shipment{}, ErrQuoteCannotTransition
	}

	return Shipment{
		ID:      fmt.Sprintf("shipment-%03d", id),
		OrderID: order.ID,
		Status:  ShipmentStatusCreated,
		Lines:   shipmentLines,
	}, nil
}
