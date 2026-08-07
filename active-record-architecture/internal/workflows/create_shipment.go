package workflows

import "active-record-architecture/internal/records"

// CreateShipment is the full-shipment convenience workflow. It delegates to
// CreatePartialShipment with no explicit line selection.
func CreateShipment(db *records.Database, orderID string, shippedBy string) (*records.Shipment, error) {
	return CreatePartialShipment(db, orderID, shippedBy, nil)
}

// CreatePartialShipment loads an order, invokes its selected-quantity
// behavior, and persists the changed order.
func CreatePartialShipment(db *records.Database, orderID string, shippedBy string, lines []records.ShipmentLine) (*records.Shipment, error) {
	if db == nil {
		return nil, records.ErrDatabaseRequired
	}

	order, err := records.FindOrder(db, orderID)
	if err != nil {
		return nil, err
	}
	shipment, err := order.CreatePartialShipment(shippedBy, lines)
	if err != nil {
		return nil, err
	}
	if err := order.Save(); err != nil {
		return nil, err
	}
	return shipment, nil
}
