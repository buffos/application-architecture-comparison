package workflows

import "active-record-architecture/internal/records"

// CreateShipment loads an order, invokes its full-shipment behavior, and
// persists the changed order. Shipment and stock records save themselves from
// the model operation.
func CreateShipment(db *records.Database, orderID string, shippedBy string) (*records.Shipment, error) {
	if db == nil {
		return nil, records.ErrDatabaseRequired
	}

	order, err := records.FindOrder(db, orderID)
	if err != nil {
		return nil, err
	}
	shipment, err := order.CreateShipment(shippedBy)
	if err != nil {
		return nil, err
	}
	if err := order.Save(); err != nil {
		return nil, err
	}
	return shipment, nil
}
