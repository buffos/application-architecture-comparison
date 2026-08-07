package workflows

import "active-record-architecture/internal/records"

// SubmitQuoteForApproval loads a Quote Active Record, invokes its lifecycle
// behavior, and persists the changed record.
func SubmitQuoteForApproval(db *records.Database, quoteID string) (*records.Quote, error) {
	if db == nil {
		return nil, records.ErrDatabaseRequired
	}

	quote, err := records.FindQuote(db, quoteID)
	if err != nil {
		return nil, err
	}
	if err := quote.SubmitForApproval(); err != nil {
		return nil, err
	}
	if err := quote.Save(); err != nil {
		return nil, err
	}
	return quote, nil
}
