package entities

import (
	"errors"
	"fmt"
	"sync/atomic"
)

const QuoteStatusDraft = "Draft"

var quoteSequence uint64

var ErrCustomerIDRequired = errors.New("customer id is required")
var ErrQuoteNotFound = errors.New("quote not found")
var ErrQuoteNotEditable = errors.New("quote is not editable")
var ErrQuoteLineQuantityInvalid = errors.New("quote line quantity must be positive")

type QuoteLine struct {
	SKU           string
	ProductName   string
	Quantity      int
	UnitPrice     int
	LineTotal     int
}

type Quote struct {
	ID         string
	CustomerID string
	Status     string
	Lines      []QuoteLine
}

func NewDraftQuote(customerID string) (Quote, error) {
	if customerID == "" {
		return Quote{}, ErrCustomerIDRequired
	}

	id := atomic.AddUint64(&quoteSequence, 1)

	return Quote{
		ID:         fmt.Sprintf("quote-%03d", id),
		CustomerID: customerID,
		Status:     QuoteStatusDraft,
	}, nil
}

func (q *Quote) AddLine(product Product, quantity int) error {
	if q.Status != QuoteStatusDraft {
		return ErrQuoteNotEditable
	}

	if quantity <= 0 {
		return ErrQuoteLineQuantityInvalid
	}

	q.Lines = append(q.Lines, QuoteLine{
		SKU:         product.SKU,
		ProductName: product.Name,
		Quantity:    quantity,
		UnitPrice:   product.BasePrice,
		LineTotal:   product.BasePrice * quantity,
	})

	return nil
}
