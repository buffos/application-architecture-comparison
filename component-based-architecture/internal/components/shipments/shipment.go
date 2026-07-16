package shipments

import "time"

type Shipment struct {
	ID         string
	OrderID    string
	CustomerID string
	Lines      []ShipmentLine
	ShippedAt  time.Time
}

type ShipmentLine struct {
	ProductSKU  string
	ProductName string
	Quantity    int
}

type ShipmentRequest struct {
	OrderID    string
	CustomerID string
	Lines      []ShipmentLine
}
