package workflows

import "active-record-architecture/internal/records"

// ApprovePaymentReview loads an order, resolves its manual payment review,
// and persists the changed order. The payment Active Record saves itself.
func ApprovePaymentReview(db *records.Database, orderID string, reviewedBy string, decision string, comment string) (*records.Order, error) {
	if db == nil {
		return nil, records.ErrDatabaseRequired
	}

	order, err := records.FindOrder(db, orderID)
	if err != nil {
		return nil, err
	}
	if err := order.ResolvePaymentReview(reviewedBy, decision, comment); err != nil {
		return nil, err
	}
	if err := order.Save(); err != nil {
		return nil, err
	}
	return order, nil
}
