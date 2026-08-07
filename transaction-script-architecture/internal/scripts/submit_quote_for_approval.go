package scripts

import (
	"errors"

	"transaction-script-architecture/internal/data"
)

var (
	ErrQuoteNotSubmittable = errors.New("quote must be in draft status")
	ErrQuoteHasNoLines     = errors.New("quote must have at least one line")
)

// SubmitQuoteForApproval owns the complete quote-submission workflow. The
// approval decision is intentionally procedural in this architecture.
func SubmitQuoteForApproval(store *data.Store, quoteID string) (data.Quote, error) {
	if store == nil {
		return data.Quote{}, ErrStoreRequired
	}

	if quoteID == "" {
		return data.Quote{}, ErrQuoteIDRequired
	}

	quote, ok := store.Quotes[quoteID]
	if !ok {
		return data.Quote{}, ErrQuoteNotFound
	}

	if quote.Status != data.QuoteStatusDraft {
		return data.Quote{}, ErrQuoteNotSubmittable
	}

	if len(quote.Lines) == 0 {
		return data.Quote{}, ErrQuoteHasNoLines
	}

	decision := EvaluateQuoteApproval(quote)
	quote.Status = data.QuoteStatusApproved
	if decision.Required {
		quote.Status = data.QuoteStatusPendingApproval
	}

	store.Quotes[quote.ID] = quote

	return quote, nil
}
