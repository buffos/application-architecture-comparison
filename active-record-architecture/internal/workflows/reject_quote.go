package workflows

import "active-record-architecture/internal/records"

// RejectQuote loads a pending Quote Active Record, applies its rejection
// behavior, and persists the resulting row.
func RejectQuote(db *records.Database, quoteID string, reviewedBy string, decisionComment string) (*records.Quote, error) {
	if db == nil {
		return nil, records.ErrDatabaseRequired
	}

	quote, err := records.FindQuote(db, quoteID)
	if err != nil {
		return nil, err
	}
	if err := quote.Reject(reviewedBy, decisionComment); err != nil {
		return nil, err
	}
	if err := quote.Save(); err != nil {
		return nil, err
	}
	return quote, nil
}
