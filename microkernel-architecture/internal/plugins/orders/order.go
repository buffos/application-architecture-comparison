package orders

import (
	"errors"
	"fmt"
	"sync/atomic"
)

var ErrOrderNotFound = errors.New("order not found")
var ErrOrderNotPayable = errors.New("order is not payable")
var ErrOrderNotShippable = errors.New("order is not shippable")
var ErrOrderNotCancellable = errors.New("order is not cancellable")

const OrderStatusPendingPayment = "PendingPayment"
const OrderStatusPaid = "Paid"
const OrderStatusShipped = "Shipped"
const OrderStatusCancelled = "Cancelled"

var orderSequence uint64

type Order struct {
	ID         string
	QuoteID    string
	CustomerID string
	Status     string
	Lines      []OrderLine
}

type OrderLine struct {
	ProductSKU      string
	ProductName     string
	ProductCategory string
	Quantity        int
	UnitPrice       int
}

func NewOrderFromApprovedQuote(quoteID string, customerID string, lines []OrderLine) Order {
	id := atomic.AddUint64(&orderSequence, 1)

	return Order{
		ID:         fmt.Sprintf("order-%03d", id),
		QuoteID:    quoteID,
		CustomerID: customerID,
		Status:     OrderStatusPendingPayment,
		Lines:      lines,
	}
}

func (o Order) TotalAmount() int {
	total := 0
	for _, line := range o.Lines {
		total += line.Quantity * line.UnitPrice
	}

	return total
}

func (o *Order) MarkPaid() error {
	if o.Status != OrderStatusPendingPayment {
		return ErrOrderNotPayable
	}

	o.Status = OrderStatusPaid
	return nil
}

func (o *Order) MarkShipped() error {
	if o.Status != OrderStatusPaid {
		return ErrOrderNotShippable
	}

	o.Status = OrderStatusShipped
	return nil
}

func (o *Order) Cancel() error {
	if o.Status == OrderStatusShipped || o.Status == OrderStatusCancelled {
		return ErrOrderNotCancellable
	}

	o.Status = OrderStatusCancelled
	return nil
}
