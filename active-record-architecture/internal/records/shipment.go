package records

import (
	"errors"
	"time"
)

const ShipmentStatusShipped = "Shipped"

var (
	ErrShipmentIDRequired = errors.New("shipment id is required")
	ErrShipmentNotFound   = errors.New("shipment not found")
)

// ShipmentLine is a shipped-quantity snapshot embedded in a Shipment Active
// Record.
type ShipmentLine struct {
	OrderLineID string
	SKU         string
	Quantity    int
}

// Shipment is an Active Record for physical fulfillment.
type Shipment struct {
	db *Database

	ID        string
	OrderID   string
	Status    string
	ShippedBy string
	ShippedAt time.Time
	Lines     []ShipmentLine
}

// FindShipment loads a Shipment Active Record from the shipment table.
func FindShipment(db *Database, id string) (*Shipment, error) {
	if db == nil {
		return nil, ErrDatabaseRequired
	}
	if id == "" {
		return nil, ErrShipmentIDRequired
	}
	row, ok := db.shipments[id]
	if !ok {
		return nil, ErrShipmentNotFound
	}
	return &Shipment{
		db:        db,
		ID:        row.ID,
		OrderID:   row.OrderID,
		Status:    row.Status,
		ShippedBy: row.ShippedBy,
		ShippedAt: row.ShippedAt,
		Lines:     cloneShipmentLines(row.Lines),
	}, nil
}

// Save writes the current Shipment Active Record to its table.
func (shipment *Shipment) Save() error {
	if shipment == nil || shipment.db == nil {
		return ErrDatabaseRequired
	}
	if shipment.ID == "" {
		return ErrShipmentIDRequired
	}
	shipment.db.shipments[shipment.ID] = shipmentRow{
		ID:        shipment.ID,
		OrderID:   shipment.OrderID,
		Status:    shipment.Status,
		ShippedBy: shipment.ShippedBy,
		ShippedAt: shipment.ShippedAt,
		Lines:     cloneShipmentLines(shipment.Lines),
	}
	return nil
}

func cloneShipmentLines(lines []ShipmentLine) []ShipmentLine {
	if lines == nil {
		return nil
	}
	clone := make([]ShipmentLine, len(lines))
	copy(clone, lines)
	return clone
}
