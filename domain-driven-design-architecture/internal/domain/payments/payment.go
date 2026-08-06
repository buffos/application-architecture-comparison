package payments

import "errors"

var (
	ErrPaymentIDRequired    = errors.New("payment id is required")
	ErrOrderIDRequired      = errors.New("order id is required")
	ErrPaymentNotCapturable = errors.New("payment is not capturable")
	ErrPaymentNotReviewable = errors.New("payment is not reviewable")
)

type PaymentID string
type OrderID string

type PaymentStatus string

const (
	PaymentStatusPending  PaymentStatus = "Pending"
	PaymentStatusReview   PaymentStatus = "Review"
	PaymentStatusCaptured PaymentStatus = "Captured"
	PaymentStatusFailed   PaymentStatus = "Failed"
)

// Payment is the aggregate root for the Payments bounded context.
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

func (p Payment) ID() PaymentID         { return p.id }
func (p Payment) OrderID() OrderID      { return p.orderID }
func (p Payment) Amount() Money         { return p.amount }
func (p Payment) Status() PaymentStatus { return p.status }

func (p *Payment) Capture() error {
	if p.status != PaymentStatusPending {
		return ErrPaymentNotCapturable
	}
	p.status = PaymentStatusCaptured
	return nil
}

func (p *Payment) RequestReview() error {
	if p.status != PaymentStatusPending {
		return ErrPaymentNotReviewable
	}
	p.status = PaymentStatusReview
	return nil
}

func (p *Payment) ApproveReview() error {
	if p.status != PaymentStatusReview {
		return ErrPaymentNotReviewable
	}
	p.status = PaymentStatusCaptured
	return nil
}

func (p *Payment) RejectReview() error {
	if p.status != PaymentStatusReview {
		return ErrPaymentNotReviewable
	}
	p.status = PaymentStatusFailed
	return nil
}

func (p *Payment) Fail() error {
	if p.status != PaymentStatusPending && p.status != PaymentStatusReview {
		return ErrPaymentNotCapturable
	}
	p.status = PaymentStatusFailed
	return nil
}
