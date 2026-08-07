package records

import "errors"

var (
	ErrPaymentIDRequired = errors.New("payment id is required")
	ErrPaymentNotFound   = errors.New("payment not found")
)

// Payment is an Active Record for a simulated payment attempt.
type Payment struct {
	db *Database

	ID              string
	OrderID         string
	Amount          int
	Status          string
	ReviewedBy      string
	DecisionComment string
}

// FindPayment loads a Payment Active Record from the payment table.
func FindPayment(db *Database, id string) (*Payment, error) {
	if db == nil {
		return nil, ErrDatabaseRequired
	}
	if id == "" {
		return nil, ErrPaymentIDRequired
	}
	row, ok := db.payments[id]
	if !ok {
		return nil, ErrPaymentNotFound
	}
	return &Payment{
		db:              db,
		ID:              row.ID,
		OrderID:         row.OrderID,
		Amount:          row.Amount,
		Status:          row.Status,
		ReviewedBy:      row.ReviewedBy,
		DecisionComment: row.DecisionComment,
	}, nil
}

// Save writes the current Payment Active Record to its table.
func (payment *Payment) Save() error {
	if payment == nil || payment.db == nil {
		return ErrDatabaseRequired
	}
	if payment.ID == "" {
		return ErrPaymentIDRequired
	}
	payment.db.payments[payment.ID] = paymentRow{
		ID:              payment.ID,
		OrderID:         payment.OrderID,
		Amount:          payment.Amount,
		Status:          payment.Status,
		ReviewedBy:      payment.ReviewedBy,
		DecisionComment: payment.DecisionComment,
	}
	return nil
}
