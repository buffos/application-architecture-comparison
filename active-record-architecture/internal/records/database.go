package records

import "fmt"

// customerRow and quoteRow stand in for rows in separate database tables.
// Active Record models translate between these private rows and the public
// model values used by callers.
type customerRow struct {
	ID     string
	Active bool
}

type quoteRow struct {
	ID               string
	CustomerID       string
	Status           string
	ConvertedOrderID string
	ReviewedBy       string
	DecisionComment  string
	Lines            []QuoteLine
}

type productRow struct {
	SKU              string
	Name             string
	Category         string
	Active           bool
	UnitPrice        int
	ReturnWindowDays int
}

type orderRow struct {
	ID            string
	SourceQuoteID string
	CustomerID    string
	Status        string
	RequestedBy   string
	PaymentStatus string
	Lines         []OrderLine
	Total         int
}

// Database is the small persistence boundary used by this lesson. Its tables
// stay private so callers must use the Active Record operations instead of
// reaching into storage directly.
type Database struct {
	customers       map[string]customerRow
	quotes          map[string]quoteRow
	products        map[string]productRow
	orders          map[string]orderRow
	nextQuoteNumber int
	nextOrderNumber int
}

func NewDatabase() *Database {
	return &Database{
		customers: make(map[string]customerRow),
		quotes:    make(map[string]quoteRow),
		products:  make(map[string]productRow),
		orders:    make(map[string]orderRow),
	}
}

func (db *Database) nextQuoteID() string {
	db.nextQuoteNumber++
	return fmt.Sprintf("quote-%03d", db.nextQuoteNumber)
}

func (db *Database) nextOrderID() string {
	db.nextOrderNumber++
	return fmt.Sprintf("order-%03d", db.nextOrderNumber)
}
