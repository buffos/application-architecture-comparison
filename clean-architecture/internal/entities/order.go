package entities

import (
	"fmt"
	"sync/atomic"
)

const OrderStatusPendingPayment = "PendingPayment"

var orderSequence uint64

var ErrQuoteNotConvertible = ErrQuoteCannotTransition

type OrderLine struct {
	SKU         string
	ProductName string
	Quantity    int
	UnitPrice   int
	LineTotal   int
}

type Order struct {
	ID            string
	CustomerID    string
	SourceQuoteID string
	Status        string
	Lines         []OrderLine
}

func NewOrderFromApprovedQuote(quote Quote) (Order, error) {
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
