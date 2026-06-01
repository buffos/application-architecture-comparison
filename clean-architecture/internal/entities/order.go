package entities

import (
	"fmt"
	"sync/atomic"
	"time"
)

const OrderStatusPendingPayment = "PendingPayment"
const OrderStatusPaid = "Paid"
const OrderStatusShipped = "Shipped"
const OrderStatusCancelled = "Cancelled"

var orderSequence uint64

var ErrQuoteNotConvertible = ErrQuoteCannotTransition

type OrderLine struct {
	SKU         string
	ProductName string
	Quantity    int
	UnitPrice   int
	LineTotal   int
	ReturnWindowDays int
}

type Order struct {
	ID            string
	CustomerID    string
	SourceQuoteID string
	Status        string
	Lines         []OrderLine
	ShippedAt     *time.Time
}

func NewOrderFromApprovedQuote(quote Quote, now time.Time) (Order, error) {
	if quote.Status != QuoteStatusApproved {
		return Order{}, ErrQuoteNotConvertible
	}

	id := atomic.AddUint64(&orderSequence, 1)
	lines := make([]OrderLine, 0, len(quote.Lines))
	for _, line := range quote.Lines {
		lines = append(lines, OrderLine{
			SKU:         line.SKU,
			ProductName: line.ProductName,
			Quantity:    line.Quantity,
			UnitPrice:   line.UnitPrice,
			LineTotal:   line.LineTotal,
			ReturnWindowDays: line.ReturnWindowDays,
		})
	}

	return Order{
		ID:            fmt.Sprintf("order-%03d", id),
		CustomerID:    quote.CustomerID,
		SourceQuoteID: quote.ID,
		Status:        OrderStatusPendingPayment,
		Lines:         lines,
	}, nil
}

func (o *Order) MarkPaid() error {
	if o.Status != OrderStatusPendingPayment {
		return ErrQuoteCannotTransition
	}

	o.Status = OrderStatusPaid
	return nil
}

func (o *Order) MarkShipped() error {
	if o.Status != OrderStatusPaid {
		return ErrQuoteCannotTransition
	}

	o.Status = OrderStatusShipped
	now := time.Now()
	o.ShippedAt = &now
	return nil
}

func (o *Order) MarkShippedAt(at time.Time) error {
	if o.Status != OrderStatusPaid {
		return ErrQuoteCannotTransition
	}

	o.Status = OrderStatusShipped
	copy := at
	o.ShippedAt = &copy
	return nil
}

func (o *Order) Cancel() error {
	if o.Status == OrderStatusShipped {
		return ErrQuoteCannotTransition
	}

	if o.Status == OrderStatusCancelled {
		return ErrQuoteCannotTransition
	}

	o.Status = OrderStatusCancelled
	return nil
}
