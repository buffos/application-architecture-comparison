package workflows

import "active-record-architecture/internal/records"

// AddQuoteLine coordinates two Active Records but leaves quote mutation and
// persistence to the records themselves.
func AddQuoteLine(db *records.Database, quoteID string, sku string, quantity int) (*records.Quote, error) {
	if db == nil {
		return nil, records.ErrDatabaseRequired
	}
	if sku == "" {
		return nil, records.ErrProductSKURequired
	}
	if quantity <= 0 {
		return nil, records.ErrQuantityInvalid
	}

	quote, err := records.FindQuote(db, quoteID)
	if err != nil {
		return nil, err
	}
	product, err := records.FindProduct(db, sku)
	if err != nil {
		return nil, err
	}
	if err := quote.AddLine(product, quantity); err != nil {
		return nil, err
	}
	if err := quote.Save(); err != nil {
		return nil, err
	}

	return quote, nil
}
