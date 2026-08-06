package scripts

import (
	"errors"
	"fmt"

	"transaction-script-architecture/internal/data"
)

var (
	ErrStoreRequired    = errors.New("store is required")
	ErrCustomerIDNeeded = errors.New("customer id is required")
	ErrCustomerNotFound = errors.New("customer not found")
	ErrCustomerInactive = errors.New("customer is inactive")
)

// CreateDraftQuote is one transaction script. It owns the complete sequence
// for this use case, including validation, lookup, ID generation, and save.
func CreateDraftQuote(store *data.Store, customerID string) (data.Quote, error) {
	if store == nil {
		return data.Quote{}, ErrStoreRequired
	}

	if customerID == "" {
		return data.Quote{}, ErrCustomerIDNeeded
	}

	customer, ok := store.Customers[customerID]
	if !ok {
		return data.Quote{}, ErrCustomerNotFound
	}

	if !customer.Active {
		return data.Quote{}, ErrCustomerInactive
	}

	store.NextQuoteNumber++
	quote := data.Quote{
		ID:         fmt.Sprintf("quote-%03d", store.NextQuoteNumber),
		CustomerID: customerID,
		Status:     data.QuoteStatusDraft,
	}

	if store.Quotes == nil {
		store.Quotes = make(map[string]data.Quote)
	}
	store.Quotes[quote.ID] = quote

	return quote, nil
}
