package scripts

import (
	"errors"

	"transaction-script-architecture/internal/data"
)

var ErrQuoteNotRejectable = errors.New("quote must be pending approval")

// RejectQuote owns the rejection transaction, including the lifecycle check,
// review metadata, and persistence of the updated passive quote record.
func RejectQuote(store *data.Store, quoteID string, reviewedBy string, decisionComment string) (data.Quote, error) {
	if store == nil {
		return data.Quote{}, ErrStoreRequired
	}

	if quoteID == "" {
		return data.Quote{}, ErrQuoteIDRequired
	}

	if reviewedBy == "" {
		return data.Quote{}, ErrReviewerRequired
	}

	quote, ok := store.Quotes[quoteID]
	if !ok {
		return data.Quote{}, ErrQuoteNotFound
	}

	if quote.Status != data.QuoteStatusPendingApproval {
		return data.Quote{}, ErrQuoteNotRejectable
	}

	quote.Status = data.QuoteStatusRejected
	quote.ReviewedBy = reviewedBy
	quote.DecisionComment = decisionComment
	store.Quotes[quote.ID] = quote

	return quote, nil
}
