package orders

import (
	"errors"
	"fmt"
	"sync/atomic"
)

var ErrOrderNotFound = errors.New("order not found")

const OrderStatusPendingPayment = "PendingPayment"

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
