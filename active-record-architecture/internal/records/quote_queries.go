package records

import "sort"

// GetQuote is the named detail-query form of FindQuote.
func GetQuote(db *Database, id string) (*Quote, error) {
	return FindQuote(db, id)
}

// ListQuotes returns reconstructed quote snapshots ordered by quote ID. Empty
// status and customer filters match every value.
func ListQuotes(db *Database, status string, customerID string) ([]*Quote, error) {
	if db == nil {
		return nil, ErrDatabaseRequired
	}

	ids := make([]string, 0, len(db.quotes))
	for id, row := range db.quotes {
		if status != "" && row.Status != status {
			continue
		}
		if customerID != "" && row.CustomerID != customerID {
			continue
		}
		ids = append(ids, id)
	}
	sort.Strings(ids)

	quotes := make([]*Quote, 0, len(ids))
	for _, id := range ids {
		quote, err := FindQuote(db, id)
		if err != nil {
			return nil, err
		}
		quotes = append(quotes, quote)
	}
	return quotes, nil
}
