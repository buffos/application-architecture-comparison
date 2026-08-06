package fulfillment

import (
	"errors"

	"domain-driven-design-architecture/internal/domain/ordering"
)

var (
	ErrShipmentIDRequired      = errors.New("shipment id is required")
	ErrOrderNotPaid            = errors.New("order is not paid")
	ErrShipmentNotDispatchable = errors.New("shipment is not dispatchable")
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

// Shipment is the aggregate root for the Fulfillment bounded context.
type Shipment struct {
	id      ShipmentID
	orderID OrderID
	status  ShipmentStatus
	lines   []ShipmentLine
}

func NewShipmentFromPaidOrder(id ShipmentID, order ordering.Order) (Shipment, error) {
	if id == "" {
		return Shipment{}, ErrShipmentIDRequired
	}
	if order.Status() != ordering.OrderStatusPaid {
		return Shipment{}, ErrOrderNotPaid
	}
	lines := make([]ShipmentLine, 0, len(order.Lines()))
	for _, line := range order.Lines() {
		lines = append(lines, ShipmentLine{sku: ProductSKU(line.ProductSKU()), quantity: line.Quantity()})
	}
	return Shipment{id: id, orderID: OrderID(order.ID()), status: ShipmentStatusPending, lines: lines}, nil
}

func (s Shipment) ID() ShipmentID         { return s.id }
func (s Shipment) OrderID() OrderID       { return s.orderID }
func (s Shipment) Status() ShipmentStatus { return s.status }
func (s Shipment) Lines() []ShipmentLine  { return append([]ShipmentLine(nil), s.lines...) }

func (s *Shipment) Dispatch() error {
	if s.status != ShipmentStatusPending {
		return ErrShipmentNotDispatchable
	}
	s.status = ShipmentStatusDispatched
	return nil
}
