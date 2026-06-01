package domain

import (
	"errors"
	"fmt"
	"sync/atomic"
)

var ErrCustomerIDRequired = errors.New("customer id is required")
var ErrQuoteNotFound = errors.New("quote not found")
var ErrQuantityMustBePositive = errors.New("quantity must be positive")
var ErrQuoteNotEditable = errors.New("quote is not editable")

const QuoteStatusDraft = "Draft"

var quoteSequence uint64

type QuoteLine struct {
	ProductSKU string
	ProductName string
	Quantity   int
	UnitPrice  int
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
		return ErrQuantityMustBePositive
	}

	q.Lines = append(q.Lines, QuoteLine{
		ProductSKU: product.SKU,
		ProductName: product.Name,
		Quantity:   quantity,
		UnitPrice:  product.UnitPrice,
	})

	return nil
}

func (q Quote) TotalQuantity() int {
	total := 0
	for _, line := range q.Lines {
		total += line.Quantity
	}

	return total
}
