package payments

import "errors"

var (
	ErrPaymentIDRequired    = errors.New("payment id is required")
	ErrOrderIDRequired      = errors.New("order id is required")
	ErrPaymentNotCapturable = errors.New("payment is not capturable")
)

type PaymentID string
type OrderID string

type PaymentStatus string

const (
	PaymentStatusPending  PaymentStatus = "Pending"
	PaymentStatusCaptured PaymentStatus = "Captured"
	PaymentStatusFailed   PaymentStatus = "Failed"
)

// Payment is the aggregate root for the Payments domain.
type Payment struct {
	id      PaymentID
	orderID OrderID
	amount  Money
	status  PaymentStatus
}

func NewPayment(id PaymentID, orderID OrderID, amount Money) (Payment, error) {
	if id == "" {
		return Payment{}, ErrPaymentIDRequired
	}
	if orderID == "" {
		return Payment{}, ErrOrderIDRequired
	}
	if amount.Currency() == "" {
		return Payment{}, ErrCurrencyRequired
	}
	return Payment{id: id, orderID: orderID, amount: amount, status: PaymentStatusPending}, nil
}

func (payment Payment) ID() PaymentID         { return payment.id }
func (payment Payment) OrderID() OrderID      { return payment.orderID }
func (payment Payment) Amount() Money         { return payment.amount }
func (payment Payment) Status() PaymentStatus { return payment.status }

func (payment *Payment) Capture() error {
	if payment.status != PaymentStatusPending {
		return ErrPaymentNotCapturable
	}
	payment.status = PaymentStatusCaptured
	return nil
}

func (payment *Payment) Fail() error {
	if payment.status != PaymentStatusPending {
		return ErrPaymentNotCapturable
	}
	payment.status = PaymentStatusFailed
	return nil
}
