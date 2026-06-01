package entities

import (
	"errors"
	"fmt"
	"sync/atomic"
)

const QuoteStatusDraft = "Draft"
const QuoteStatusSubmitted = "Submitted"
const QuoteStatusPendingApproval = "PendingApproval"
const QuoteStatusApproved = "Approved"

var quoteSequence uint64

var ErrCustomerIDRequired = errors.New("customer id is required")
var ErrQuoteNotFound = errors.New("quote not found")
var ErrQuoteNotEditable = errors.New("quote is not editable")
var ErrQuoteLineQuantityInvalid = errors.New("quote line quantity must be positive")
var ErrQuoteCannotSubmitWithoutLines = errors.New("quote must have at least one line before submission")
var ErrQuoteCannotTransition = errors.New("quote cannot transition from its current status")

type QuoteLine struct {
	SKU           string
	ProductName   string
	ProductCategory string
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
		SKU:             product.SKU,
		ProductName:     product.Name,
		ProductCategory: product.Category,
		Quantity:        quantity,
		UnitPrice:       product.BasePrice,
		LineTotal:       product.BasePrice * quantity,
	})

	return nil
}

func (q *Quote) Submit(requiresApproval bool) error {
	if q.Status != QuoteStatusDraft {
		return ErrQuoteCannotTransition
	}

	if len(q.Lines) == 0 {
		return ErrQuoteCannotSubmitWithoutLines
	}

	if requiresApproval {
		q.Status = QuoteStatusPendingApproval
		return nil
	}

	q.Status = QuoteStatusApproved
	return nil
}
