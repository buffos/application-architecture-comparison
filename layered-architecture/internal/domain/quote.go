package domain

import (
	"errors"
	"fmt"
	"sync/atomic"
)

const QuoteStatusDraft = "Draft"
const QuoteStatusPendingApproval = "PendingApproval"
const QuoteStatusApproved = "Approved"
const QuoteStatusRejected = "Rejected"

var quoteSequence uint64

var ErrQuoteNotFound = errors.New("quote not found")
var ErrQuoteNotEditable = errors.New("quote is no longer editable")
var ErrQuoteCannotTransition = errors.New("quote cannot transition from its current status")
var ErrQuoteLineProductRequired = errors.New("product name is required")
var ErrQuoteLineQuantityInvalid = errors.New("quantity must be positive")
var ErrQuoteCannotSubmitWithoutLines = errors.New("quote must have at least one line before submission")
var ErrQuoteNotApproved = errors.New("quote must be approved before conversion")

type QuoteLine struct {
	SKU                 string
	ProductCategory     string
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

func (q *Quote) AddLine(sku string, productCategory string, productName string, quantity int) error {
	if q.Status != QuoteStatusDraft {
		return ErrQuoteNotEditable
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
		ProductCategory:     productCategory,
		ProductNameSnapshot: productName,
		Quantity:            quantity,
	})

	return nil
}

func (q *Quote) Submit() error {
	if q.Status != QuoteStatusDraft {
		return ErrQuoteCannotTransition
	}

	if len(q.Lines) == 0 {
		return ErrQuoteCannotSubmitWithoutLines
	}

	if q.RequiresApproval() {
		q.Status = QuoteStatusPendingApproval
		return nil
	}

	q.Status = QuoteStatusApproved

	return nil
}

func (q Quote) RequiresApproval() bool {
	for _, line := range q.Lines {
		if line.ProductCategory == "CustomBuild" {
			return true
		}
	}

	return false
}

func (q *Quote) Approve() error {
	if q.Status != QuoteStatusPendingApproval {
		return ErrQuoteCannotTransition
	}

	q.Status = QuoteStatusApproved
	return nil
}

func (q *Quote) Reject() error {
	if q.Status != QuoteStatusPendingApproval {
		return ErrQuoteCannotTransition
	}

	q.Status = QuoteStatusRejected
	return nil
}
