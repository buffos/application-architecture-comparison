package orders

import (
	"errors"
	"fmt"
	"sync/atomic"

	"modular-monolith/internal/modules/quotes"
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

func NewOrderFromApprovedQuote(quote quotes.ApprovedQuote) Order {
	id := atomic.AddUint64(&orderSequence, 1)

	lines := make([]OrderLine, 0, len(quote.Lines))
	for _, line := range quote.Lines {
		lines = append(lines, OrderLine{
			ProductSKU:      line.ProductSKU,
			ProductName:     line.ProductName,
			ProductCategory: line.ProductCategory,
			Quantity:        line.Quantity,
			UnitPrice:       line.UnitPrice,
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
