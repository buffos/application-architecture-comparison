package workflows

import "active-record-architecture/internal/records"

// CancelOrder loads an order, invokes its cancellation behavior, and saves
// the changed order. Stock rows are saved by Order.Cancel.
func CancelOrder(db *records.Database, orderID string, cancelledBy string, reason string) (*records.Order, error) {
	if db == nil {
		return nil, records.ErrDatabaseRequired
	}

	order, err := records.FindOrder(db, orderID)
	if err != nil {
		return nil, err
	}
	if err := order.Cancel(cancelledBy, reason); err != nil {
		return nil, err
	}
	if err := order.Save(); err != nil {
		return nil, err
	}
	return order, nil
}
