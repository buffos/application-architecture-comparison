package fulfillment

import (
	"errors"

	"rich-domain-model-architecture/internal/domain/ordering"
)

var (
	ErrShipmentIDRequired       = errors.New("shipment id is required")
	ErrOrderNotPaid             = errors.New("order is not paid")
	ErrShipmentNotDispatchable  = errors.New("shipment is not dispatchable")
	ErrShipmentSelectionInvalid = errors.New("shipment selection is invalid")
)

type ShipmentID string
type OrderID string
type ProductSKU string

type ShipmentStatus string

const (
	ShipmentStatusPending    ShipmentStatus = "Pending"
	ShipmentStatusDispatched ShipmentStatus = "Dispatched"
)

type ShipmentLine struct {
	sku      ProductSKU
	quantity int
}

func (line ShipmentLine) ProductSKU() ProductSKU { return line.sku }
func (line ShipmentLine) Quantity() int          { return line.quantity }

// Shipment is the aggregate root for the Fulfillment domain.
type Shipment struct {
	id      ShipmentID
	orderID OrderID
	status  ShipmentStatus
	lines   []ShipmentLine
}

func NewShipmentFromPaidOrder(id ShipmentID, order ordering.Order) (Shipment, error) {
	return NewShipmentFromOrderSelection(id, order, nil)
}

func NewShipmentFromOrderSelection(id ShipmentID, order ordering.Order, selections []ordering.ShipmentSelection) (Shipment, error) {
	if id == "" {
		return Shipment{}, ErrShipmentIDRequired
	}
	if order.Status() != ordering.OrderStatusPaid && order.Status() != ordering.OrderStatusPartiallyShipped {
		return Shipment{}, ErrOrderNotPaid
	}
	if len(selections) == 0 {
		for _, line := range order.Lines() {
			if remaining := line.Quantity() - line.ShippedQuantity(); remaining > 0 {
				selections = append(selections, ordering.ShipmentSelection{
					ProductSKU: ordering.ProductSKU(line.ProductSKU()),
					Quantity:   remaining,
				})
			}
		}
	}

	requested := make(map[ordering.ProductSKU]int, len(selections))
	for _, selection := range selections {
		if selection.Quantity <= 0 {
			return Shipment{}, ErrShipmentSelectionInvalid
		}
		requested[selection.ProductSKU] += selection.Quantity
	}
	for sku, quantity := range requested {
		matched := false
		for _, line := range order.Lines() {
			if ordering.ProductSKU(line.ProductSKU()) != sku {
				continue
			}
			if quantity > line.Quantity()-line.ShippedQuantity() {
				return Shipment{}, ErrShipmentSelectionInvalid
			}
			matched = true
			break
		}
		if !matched {
			return Shipment{}, ErrShipmentSelectionInvalid
		}
	}

	lines := make([]ShipmentLine, 0, len(requested))
	for sku, quantity := range requested {
		lines = append(lines, ShipmentLine{sku: ProductSKU(sku), quantity: quantity})
	}
	return Shipment{id: id, orderID: OrderID(order.ID()), status: ShipmentStatusPending, lines: lines}, nil
}

func (shipment Shipment) ID() ShipmentID         { return shipment.id }
func (shipment Shipment) OrderID() OrderID       { return shipment.orderID }
func (shipment Shipment) Status() ShipmentStatus { return shipment.status }
func (shipment Shipment) Lines() []ShipmentLine {
	return append([]ShipmentLine(nil), shipment.lines...)
}

func (shipment *Shipment) Dispatch() error {
	if shipment.status != ShipmentStatusPending {
		return ErrShipmentNotDispatchable
	}
	shipment.status = ShipmentStatusDispatched
	return nil
}
