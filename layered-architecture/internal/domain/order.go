package domain

import (
	"errors"
	"fmt"
	"sync/atomic"
)

const OrderStatusReadyForPayment = "ReadyForPayment"

var orderSequence uint64

var ErrOrderNotFound = errors.New("order not found")
var ErrQuoteNotConvertible = errors.New("quote must be approved before conversion")

type OrderLine struct {
	SKU                 string
	ProductNameSnapshot string
	Quantity            int
}

type Order struct {
	ID            string
	SourceQuoteID string
	CustomerID    string
	Status        string
	Lines         []OrderLine
}

func NewOrderFromQuote(quote Quote) (Order, error) {
	if quote.Status != QuoteStatusApproved {
		return Order{}, ErrQuoteNotConvertible
	}

	id := atomic.AddUint64(&orderSequence, 1)
	lines := make([]OrderLine, 0, len(quote.Lines))

	for _, line := range quote.Lines {
		lines = append(lines, OrderLine{
			SKU:                 line.SKU,
			ProductNameSnapshot: line.ProductNameSnapshot,
			Quantity:            line.Quantity,
		})
	}

	return Order{
		ID:            fmt.Sprintf("order-%03d", id),
		SourceQuoteID: quote.ID,
		CustomerID:    quote.CustomerID,
		Status:        OrderStatusReadyForPayment,
		Lines:         lines,
	}, nil
}
