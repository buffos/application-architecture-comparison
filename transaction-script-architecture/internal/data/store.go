package data

const (
	QuoteStatusDraft           = "Draft"
	QuoteStatusPendingApproval = "PendingApproval"
	QuoteStatusApproved        = "Approved"
	QuoteStatusRejected        = "Rejected"
)

// Customer is a passive data record used by transaction scripts.
type Customer struct {
	ID     string
	Active bool
}

// Product is a passive data record used by transaction scripts.
type Product struct {
	SKU       string
	Name      string
	Category  string
	Active    bool
	UnitPrice int
}

// QuoteLine is a passive data record embedded in a Quote.
type QuoteLine struct {
	ProductCategory     string
	SKU                 string
	ProductNameSnapshot string
	Quantity            int
	UnitPrice           int
	LineTotal           int
}

// Quote is a passive data record used by transaction scripts.
type Quote struct {
	ID              string
	CustomerID      string
	Status          string
	ReviewedBy      string
	DecisionComment string
	Lines           []QuoteLine
}

// Store exposes the in-memory data shape directly so scripts can coordinate a
// transaction without a repository or domain-model abstraction.
type Store struct {
	Customers       map[string]Customer
	Products        map[string]Product
	Quotes          map[string]Quote
	NextQuoteNumber int
}

func NewStore() *Store {
	return &Store{
		Customers: make(map[string]Customer),
		Products:  make(map[string]Product),
		Quotes:    make(map[string]Quote),
	}
}
