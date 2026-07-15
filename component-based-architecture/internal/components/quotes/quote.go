package quotes

import "errors"

var ErrCustomerIDRequired = errors.New("customer id is required")

const QuoteStatusDraft = "Draft"

type Quote struct {
	ID         string
	CustomerID string
	Status     string
}

type CreateDraftQuoteCommand struct {
	CustomerID string
}

type CreateDraftQuoteResult struct {
	QuoteID    string
	CustomerID string
	Status     string
}
