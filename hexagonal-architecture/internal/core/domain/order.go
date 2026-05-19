package domain

import (
	"errors"
	"fmt"
	"sync/atomic"
	"time"
)

const OrderStatusReadyForPayment = "ReadyForPayment"
const OrderStatusPaymentReview = "PaymentReview"
const OrderStatusReadyForFulfillment = "ReadyForFulfillment"
const OrderStatusShipped = "Shipped"
const OrderStatusCancelled = "Cancelled"

var orderSequence uint64

var ErrOrderNotFound = errors.New("order not found")
var ErrQuoteNotApproved = errors.New("quote must be approved before conversion")
var ErrPaymentFailed = errors.New("payment failed")
var ErrPaymentReviewNotAllowed = errors.New("payment review is not allowed")
var ErrOrderActorRequired = errors.New("actor is required")
var ErrOrderCancellationNotAllowed = errors.New("order cancellation is not allowed after shipment")

type OrderLine struct {
	SKU               string
	ProductName       string
	ProductCategory   string
	Quantity          int
	BaseUnitPrice     int
	AdjustedUnitPrice int
	LineTotal         int
	ReturnWindowDays  int
}

type Order struct {
	ID                string
	SourceQuoteID     string
	CustomerID        string
	Status            string
	PaymentStatus     string
	PaymentReviewedBy string
	ShippedAt         time.Time
	Lines             []OrderLine
}

func NewOrderFromQuote(quote Quote) (Order, error) {
	if quote.Status != QuoteStatusApproved {
		return Order{}, ErrQuoteNotApproved
	}

	id := atomic.AddUint64(&orderSequence, 1)
	lines := make([]OrderLine, 0, len(quote.Lines))

	for _, line := range quote.Lines {
		lines = append(lines, OrderLine{
			SKU:               line.SKU,
			ProductName:       line.ProductName,
			ProductCategory:   line.ProductCategory,
			Quantity:          line.Quantity,
			BaseUnitPrice:     line.BaseUnitPrice,
			AdjustedUnitPrice: line.AdjustedUnitPrice,
			LineTotal:         line.LineTotal,
			ReturnWindowDays:  line.ReturnWindowDays,
		})
	}

	return Order{
		ID:            fmt.Sprintf("order-%03d", id),
		SourceQuoteID: quote.ID,
		CustomerID:    quote.CustomerID,
		Status:        OrderStatusReadyForPayment,
		PaymentStatus: "Pending",
		Lines:         lines,
	}, nil
}

func (o *Order) AcceptPayment() {
	o.PaymentStatus = "Accepted"
	o.Status = OrderStatusReadyForFulfillment
}

func (o *Order) MarkPaymentReview() {
	o.PaymentStatus = "ManualReview"
	o.Status = OrderStatusPaymentReview
}

func (o *Order) FailPayment() {
	o.PaymentStatus = "Failed"
}

func (o *Order) ApprovePaymentReview(reviewedBy string) error {
	if o.Status != OrderStatusPaymentReview {
		return ErrPaymentReviewNotAllowed
	}
	if reviewedBy == "" {
		return ErrOrderActorRequired
	}

	o.PaymentReviewedBy = reviewedBy
	o.AcceptPayment()
	return nil
}

func (o *Order) MarkShipped(shippedAt time.Time) error {
	if o.Status != OrderStatusReadyForFulfillment {
		return ErrShipmentNotAllowedUntilPaymentAccepted
	}

	o.Status = OrderStatusShipped
	o.ShippedAt = shippedAt
	return nil
}

func (o *Order) Cancel() error {
	if o.Status == OrderStatusShipped {
		return ErrOrderCancellationNotAllowed
	}

	o.Status = OrderStatusCancelled
	return nil
}
