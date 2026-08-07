package workflows

import "active-record-architecture/internal/records"

// ConvertQuoteToOrder loads an approved Quote Active Record, asks it to
// create an independent Order snapshot, and saves both records.
func ConvertQuoteToOrder(db *records.Database, quoteID string, requestedBy string) (*records.Order, error) {
	if db == nil {
		return nil, records.ErrDatabaseRequired
	}

	quote, err := records.FindQuote(db, quoteID)
	if err != nil {
		return nil, err
	}
	order, err := quote.ConvertToOrder(requestedBy)
	if err != nil {
		return nil, err
	}
	if err := order.ReserveStock(); err != nil {
		return nil, err
	}
	if err := order.Save(); err != nil {
		return nil, err
	}
	if err := quote.Save(); err != nil {
		return nil, err
	}
	return order, nil
}
