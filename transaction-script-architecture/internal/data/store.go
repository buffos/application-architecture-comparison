package data

const (
	QuoteStatusDraft           = "Draft"
	QuoteStatusPendingApproval = "PendingApproval"
	QuoteStatusApproved        = "Approved"
	QuoteStatusRejected        = "Rejected"
	QuoteStatusConverted       = "Converted"
)

// Customer is a passive data record used by transaction scripts.
type Customer struct {
	ID     string
	Active bool
}

// Product is a passive data record used by transaction scripts.
type Product struct {
	SKU                 string
	Name                string
	Category            string
	Active              bool
	UnitPrice           int
	StockShortagePolicy string
}

// QuoteLine is a passive data record embedded in a Quote.
type QuoteLine struct {
	ProductCategory     string
	SKU                 string
	ProductNameSnapshot string
	Quantity            int
	UnitPrice           int
	DiscountAmount      int
	LineTotal           int
}

// Quote is a passive data record used by transaction scripts.
type Quote struct {
	ID               string
	CustomerID       string
	Status           string
	ConvertedOrderID string
	ReviewedBy       string
	DecisionComment  string
	Lines            []QuoteLine
}

// Store exposes the in-memory data shape directly so scripts can coordinate a
// transaction without a repository or domain-model abstraction.
type Store struct {
	Customers         map[string]Customer
	Products          map[string]Product
	Stocks            map[string]StockRecord
	Quotes            map[string]Quote
	Orders            map[string]Order
	Payments          map[string]Payment
	NextQuoteNumber   int
	NextOrderNumber   int
	NextPaymentNumber int
}

func NewStore() *Store {
	return &Store{
		Customers: make(map[string]Customer),
		Products:  make(map[string]Product),
		Stocks:    make(map[string]StockRecord),
		Quotes:    make(map[string]Quote),
		Orders:    make(map[string]Order),
		Payments:  make(map[string]Payment),
	}
}
