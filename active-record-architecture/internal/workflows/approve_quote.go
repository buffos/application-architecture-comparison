package workflows

import "active-record-architecture/internal/records"

// ApproveQuote loads a pending Quote Active Record, applies its approval
// behavior, and saves the resulting row.
func ApproveQuote(db *records.Database, quoteID string, reviewedBy string, decisionComment string) (*records.Quote, error) {
	if db == nil {
		return nil, records.ErrDatabaseRequired
	}

	quote, err := records.FindQuote(db, quoteID)
	if err != nil {
		return nil, err
	}
	if err := quote.Approve(reviewedBy, decisionComment); err != nil {
		return nil, err
	}
	if err := quote.Save(); err != nil {
		return nil, err
	}
	return quote, nil
}
