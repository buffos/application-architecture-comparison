package data

const QuoteStatusDraft = "Draft"

// Customer is a passive data record used by transaction scripts.
type Customer struct {
	ID     string
	Active bool
}

// Quote is a passive data record used by transaction scripts.
type Quote struct {
	ID         string
	CustomerID string
	Status     string
}

// Store exposes the in-memory data shape directly so scripts can coordinate a
// transaction without a repository or domain-model abstraction.
type Store struct {
	Customers       map[string]Customer
	Quotes          map[string]Quote
	NextQuoteNumber int
}

func NewStore() *Store {
	return &Store{
		Customers: make(map[string]Customer),
		Quotes:    make(map[string]Quote),
	}
}
