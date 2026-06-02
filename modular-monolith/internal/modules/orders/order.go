package orders

import (
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"modular-monolith/internal/modules/quotes"
)

var ErrOrderNotFound = errors.New("order not found")
var ErrOrderNotPayable = errors.New("order is not payable")
var ErrOrderNotInPaymentReview = errors.New("order is not in payment review")
var ErrOrderNotShippable = errors.New("order is not shippable")
var ErrOrderNotCancellable = errors.New("order is not cancellable")
var ErrOrderNotReturnable = errors.New("order is not returnable")

const (
	OrderStatusPendingPayment = "PendingPayment"
	OrderStatusPaymentReview  = "PaymentReview"
	OrderStatusPaid           = "Paid"
	OrderStatusShipped        = "Shipped"
	OrderStatusCancelled      = "Cancelled"
)

var orderSequence uint64

type Order struct {
	ID         string
	QuoteID    string
	CustomerID string
	Status     string
	Lines      []OrderLine
	ShippedAt  time.Time
}

type OrderLine struct {
	ProductSKU       string
	ProductName      string
	ProductCategory  string
	Quantity         int
	UnitPrice        int
	ReturnWindowDays int
}

func (o *Order) MarkPaid() error {
	if o.Status != OrderStatusPendingPayment && o.Status != OrderStatusPaymentReview {
		return ErrOrderNotPayable
	}

	o.Status = OrderStatusPaid
	return nil
}

func (o *Order) MarkPaymentReview() error {
	if o.Status != OrderStatusPendingPayment {
		return ErrOrderNotPayable
	}

	o.Status = OrderStatusPaymentReview
	return nil
}

func (o *Order) ApprovePaymentReview() error {
	if o.Status != OrderStatusPaymentReview {
		return ErrOrderNotInPaymentReview
	}

	o.Status = OrderStatusPaid
	return nil
}

func (o *Order) MarkShipped(shippedAt time.Time) error {
	if o.Status != OrderStatusPaid {
		return ErrOrderNotShippable
	}

	o.Status = OrderStatusShipped
	o.ShippedAt = shippedAt
	return nil
}

func (o *Order) Cancel() error {
	if o.Status == OrderStatusShipped || o.Status == OrderStatusCancelled {
		return ErrOrderNotCancellable
	}

	o.Status = OrderStatusCancelled
	return nil
}

func (o Order) EnsureReturnable() error {
	if o.Status != OrderStatusShipped {
		return ErrOrderNotReturnable
	}

	return nil
}

func NewOrderFromApprovedQuote(quote quotes.ApprovedQuote) Order {
	id := atomic.AddUint64(&orderSequence, 1)

	lines := make([]OrderLine, 0, len(quote.Lines))
	for _, line := range quote.Lines {
		lines = append(lines, OrderLine{
			ProductSKU:       line.ProductSKU,
			ProductName:      line.ProductName,
			ProductCategory:  line.ProductCategory,
			Quantity:         line.Quantity,
			UnitPrice:        line.UnitPrice,
			ReturnWindowDays: line.ReturnWindowDays,
		})
	}

	return Order{
		ID:         fmt.Sprintf("order-%03d", id),
		QuoteID:    quote.QuoteID,
		CustomerID: quote.CustomerID,
		Status:     OrderStatusPendingPayment,
		Lines:      lines,
	}
}
