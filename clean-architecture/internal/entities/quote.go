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

type Quote struct {
	ID         string
	CustomerID string
	Status     string
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
