package domain

import (
	"errors"
	"fmt"
	"sync/atomic"
)

var ErrShipmentNotFound = errors.New("shipment not found")
var ErrShipmentQuantityMustBePositive = errors.New("shipment quantity must be positive")
var ErrShipmentLineNotFound = errors.New("shipment line not found")
var ErrShipmentQuantityExceedsRemaining = errors.New("shipment quantity exceeds remaining quantity")
var ErrShipmentHasNoRemainingQuantity = errors.New("shipment has no remaining quantity")

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

func NewShipmentFromOrder(order Order, lines []ShipmentLine) (Shipment, error) {
	if err := order.EnsureShippable(); err != nil {
		return Shipment{}, err
	}

	resolved, err := resolveShipmentLines(order, lines)
	if err != nil {
		return Shipment{}, err
	}

	id := atomic.AddUint64(&shipmentSequence, 1)

	return Shipment{
		ID:      fmt.Sprintf("shipment-%03d", id),
		OrderID: order.ID,
		Lines:   resolved,
	}, nil
}

func resolveShipmentLines(order Order, requested []ShipmentLine) ([]ShipmentLine, error) {
	if err := order.EnsureShippable(); err != nil {
		return nil, err
	}

	remainingBySKU := make(map[string]int, len(order.Lines))
	for _, line := range order.Lines {
		remaining := line.Quantity - line.ShippedQuantity
		if remaining > 0 {
			remainingBySKU[line.ProductSKU] = remaining
		}
	}

	if len(requested) == 0 {
		lines := make([]ShipmentLine, 0, len(remainingBySKU))
		for _, line := range order.Lines {
			remaining := remainingBySKU[line.ProductSKU]
			if remaining == 0 {
				continue
			}

			lines = append(lines, ShipmentLine{
				ProductSKU: line.ProductSKU,
				Quantity:   remaining,
			})
		}

		if len(lines) == 0 {
			return nil, ErrShipmentHasNoRemainingQuantity
		}

		return lines, nil
	}

	accumulated := make(map[string]int)
	lines := make([]ShipmentLine, 0, len(requested))
	for _, line := range requested {
		if line.Quantity <= 0 {
			return nil, ErrShipmentQuantityMustBePositive
		}

		remaining, ok := remainingBySKU[line.ProductSKU]
		if !ok {
			return nil, ErrShipmentLineNotFound
		}

		accumulated[line.ProductSKU] += line.Quantity
		if accumulated[line.ProductSKU] > remaining {
			return nil, ErrShipmentQuantityExceedsRemaining
		}

		lines = append(lines, ShipmentLine{
			ProductSKU: line.ProductSKU,
			Quantity:   line.Quantity,
		})
	}

	return lines, nil
}
