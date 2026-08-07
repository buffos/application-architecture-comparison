package records

import "sort"

// GetShipment is the named detail-query form of FindShipment.
func GetShipment(db *Database, id string) (*Shipment, error) {
	return FindShipment(db, id)
}

// ListShipments returns reconstructed shipment snapshots ordered by shipment
// ID. An empty order ID lists every shipment.
func ListShipments(db *Database, orderID string) ([]*Shipment, error) {
	if db == nil {
		return nil, ErrDatabaseRequired
	}

	ids := make([]string, 0, len(db.shipments))
	for id, row := range db.shipments {
		if orderID != "" && row.OrderID != orderID {
			continue
		}
		ids = append(ids, id)
	}
	sort.Strings(ids)

	shipments := make([]*Shipment, 0, len(ids))
	for _, id := range ids {
		shipment, err := FindShipment(db, id)
		if err != nil {
			return nil, err
		}
		shipments = append(shipments, shipment)
	}
	return shipments, nil
}
