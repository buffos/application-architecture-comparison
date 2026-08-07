package data

const ShipmentStatusShipped = "Shipped"

// ShipmentLine is a passive shipped-quantity record.
type ShipmentLine struct {
	OrderLineID string
	SKU         string
	Quantity    int
}

// Shipment is a passive fulfillment record.
type Shipment struct {
	ID        string
	OrderID   string
	Status    string
	ShippedBy string
	Lines     []ShipmentLine
}
