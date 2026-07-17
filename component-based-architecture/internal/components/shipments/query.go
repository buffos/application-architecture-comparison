package shipments

import (
	"errors"
	"time"
)

var ErrShipmentNotFound = errors.New("shipment not found")

// Reader is the public read contract provided by Shipments.
type Reader interface {
	GetShipment(query GetShipmentQuery) (ShipmentDetails, error)
	ListShipments(query ListShipmentsQuery) []ShipmentSummary
}

type GetShipmentQuery struct{ ShipmentID string }
type ListShipmentsQuery struct{ OrderID string }

type ShipmentDetails struct {
	ShipmentID string
	OrderID    string
	CustomerID string
	ShippedAt  time.Time
	LineCount  int
	Lines      []ShipmentLineDetails
}

type ShipmentLineDetails struct {
	ProductSKU  string
	ProductName string
	Quantity    int
}

type ShipmentSummary struct {
	ShipmentID string
	OrderID    string
	CustomerID string
	LineCount  int
}

func (c *Component) GetShipment(query GetShipmentQuery) (ShipmentDetails, error) {
	shipment, ok := c.shipments[query.ShipmentID]
	if !ok {
		return ShipmentDetails{}, ErrShipmentNotFound
	}
	return shipmentDetails(shipment), nil
}

func (c *Component) ListShipments(query ListShipmentsQuery) []ShipmentSummary {
	shipments := make([]ShipmentSummary, 0, len(c.shipments))
	for _, shipment := range c.shipments {
		if query.OrderID != "" && shipment.OrderID != query.OrderID {
			continue
		}
		shipments = append(shipments, ShipmentSummary{ShipmentID: shipment.ID, OrderID: shipment.OrderID, CustomerID: shipment.CustomerID, LineCount: len(shipment.Lines)})
	}
	return shipments
}

func shipmentDetails(shipment Shipment) ShipmentDetails {
	lines := make([]ShipmentLineDetails, 0, len(shipment.Lines))
	for _, line := range shipment.Lines {
		lines = append(lines, ShipmentLineDetails{ProductSKU: line.ProductSKU, ProductName: line.ProductName, Quantity: line.Quantity})
	}
	return ShipmentDetails{ShipmentID: shipment.ID, OrderID: shipment.OrderID, CustomerID: shipment.CustomerID, ShippedAt: shipment.ShippedAt, LineCount: len(lines), Lines: lines}
}
