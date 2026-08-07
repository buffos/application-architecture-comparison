package records

import "errors"

var (
	ErrRefundIDRequired = errors.New("refund id is required")
	ErrRefundNotFound   = errors.New("refund not found")
)

// Refund is a passive financial follow-up Active Record for a return request.
type Refund struct {
	db *Database

	ID              string
	ReturnRequestID string
	OrderID         string
	Amount          int
	Status          string
}

// FindRefund loads a Refund Active Record from the refunds table.
func FindRefund(db *Database, id string) (*Refund, error) {
	if db == nil {
		return nil, ErrDatabaseRequired
	}
	if id == "" {
		return nil, ErrRefundIDRequired
	}

	row, ok := db.refunds[id]
	if !ok {
		return nil, ErrRefundNotFound
	}
	return &Refund{
		db:              db,
		ID:              row.ID,
		ReturnRequestID: row.ReturnRequestID,
		OrderID:         row.OrderID,
		Amount:          row.Amount,
		Status:          row.Status,
	}, nil
}

// Save writes the current Refund Active Record to its table.
func (refund *Refund) Save() error {
	if refund == nil || refund.db == nil {
		return ErrDatabaseRequired
	}
	if refund.ID == "" {
		return ErrRefundIDRequired
	}
	refund.db.refunds[refund.ID] = refundRow{
		ID:              refund.ID,
		ReturnRequestID: refund.ReturnRequestID,
		OrderID:         refund.OrderID,
		Amount:          refund.Amount,
		Status:          refund.Status,
	}
	return nil
}
