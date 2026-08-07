package workflows

import "active-record-architecture/internal/records"

// CapturePayment loads an order, invokes its payment behavior, and persists
// the changed order. The payment Active Record is saved by the model method.
func CapturePayment(db *records.Database, orderID string, outcome string) (*records.Order, error) {
	if db == nil {
		return nil, records.ErrDatabaseRequired
	}

	order, err := records.FindOrder(db, orderID)
	if err != nil {
		return nil, err
	}
	if _, err := order.CapturePayment(outcome); err != nil {
		return nil, err
	}
	if err := order.Save(); err != nil {
		return nil, err
	}
	return order, nil
}
