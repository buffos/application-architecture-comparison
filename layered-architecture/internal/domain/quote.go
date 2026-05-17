package domain

import (
	"errors"
	"fmt"
	"sync/atomic"
)

const QuoteStatusDraft = "Draft"
const QuoteStatusSubmitted = "Submitted"

var quoteSequence uint64

var ErrQuoteNotFound = errors.New("quote not found")
var ErrQuoteAlreadySubmitted = errors.New("quote is already submitted")
var ErrQuoteLineProductRequired = errors.New("product name is required")
var ErrQuoteLineQuantityInvalid = errors.New("quantity must be positive")
var ErrQuoteCannotSubmitWithoutLines = errors.New("quote must have at least one line before submission")

type QuoteLine struct {
	SKU                 string
	ProductNameSnapshot string
	Quantity            int
}

type Quote struct {
	ID         string
	CustomerID string
	Status     string
	Lines      []QuoteLine
}

func NewDraftQuote(customerID string) (Quote, error) {
	if customerID == "" {
		return Quote{}, errors.New("customer id is required")
	}

	id := atomic.AddUint64(&quoteSequence, 1)

	return Quote{
		ID:         fmt.Sprintf("quote-%03d", id),
		CustomerID: customerID,
		Status:     QuoteStatusDraft,
	}, nil
}

func (q *Quote) AddLine(sku string, productName string, quantity int) error {
	if q.Status != QuoteStatusDraft {
		return ErrQuoteAlreadySubmitted
	}

	if sku == "" {
		return ErrProductSKURequired
	}

	if productName == "" {
		return ErrQuoteLineProductRequired
	}

	if quantity <= 0 {
		return ErrQuoteLineQuantityInvalid
	}

	q.Lines = append(q.Lines, QuoteLine{
		SKU:                 sku,
		ProductNameSnapshot: productName,
		Quantity:            quantity,
	})

	return nil
}

func (q *Quote) Submit() error {
	if q.Status != QuoteStatusDraft {
		return ErrQuoteAlreadySubmitted
	}

	if len(q.Lines) == 0 {
		return ErrQuoteCannotSubmitWithoutLines
	}

	q.Status = QuoteStatusSubmitted

	return nil
}
