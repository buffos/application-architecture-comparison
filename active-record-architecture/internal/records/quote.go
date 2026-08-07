package records

import "errors"

const QuoteStatusDraft = "Draft"

var (
	ErrQuoteIDRequired         = errors.New("quote id is required")
	ErrQuoteCustomerIDRequired = errors.New("quote customer id is required")
	ErrQuoteStatusRequired     = errors.New("quote status is required")
	ErrQuoteNotFound           = errors.New("quote not found")
)

// Quote is an Active Record. It contains the quote fields and the database
// connection used by Save and FindQuote.
type Quote struct {
	db *Database

	ID         string
	CustomerID string
	Status     string
}

// NewDraftQuote creates an unsaved draft quote for an existing active
// customer. The caller explicitly persists the returned Active Record with
// Quote.Save.
func NewDraftQuote(db *Database, customerID string) (*Quote, error) {
	if db == nil {
		return nil, ErrDatabaseRequired
	}
	if customerID == "" {
		return nil, ErrCustomerIDRequired
	}

	customer, err := FindCustomer(db, customerID)
	if err != nil {
		return nil, err
	}
	if !customer.Active {
		return nil, ErrCustomerInactive
	}

	return &Quote{
		db:         db,
		ID:         db.nextQuoteID(),
		CustomerID: customer.ID,
		Status:     QuoteStatusDraft,
	}, nil
}

// FindQuote loads a Quote Active Record from the quote table.
func FindQuote(db *Database, id string) (*Quote, error) {
	if db == nil {
		return nil, ErrDatabaseRequired
	}
	if id == "" {
		return nil, ErrQuoteIDRequired
	}

	row, ok := db.quotes[id]
	if !ok {
		return nil, ErrQuoteNotFound
	}

	return &Quote{
		db:         db,
		ID:         row.ID,
		CustomerID: row.CustomerID,
		Status:     row.Status,
	}, nil
}

// Save writes the current Quote Active Record to its table.
func (quote *Quote) Save() error {
	if quote == nil || quote.db == nil {
		return ErrDatabaseRequired
	}
	if quote.ID == "" {
		return ErrQuoteIDRequired
	}
	if quote.CustomerID == "" {
		return ErrQuoteCustomerIDRequired
	}
	if quote.Status == "" {
		return ErrQuoteStatusRequired
	}
	if _, ok := quote.db.customers[quote.CustomerID]; !ok {
		return ErrCustomerNotFound
	}

	quote.db.quotes[quote.ID] = quoteRow{
		ID:         quote.ID,
		CustomerID: quote.CustomerID,
		Status:     quote.Status,
	}
	return nil
}
