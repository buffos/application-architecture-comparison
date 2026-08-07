package records

import (
	"fmt"
	"time"
)

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
	SKU                 string
	Name                string
	Category            string
	Active              bool
	UnitPrice           int
	ReturnWindowDays    int
	StockShortagePolicy string
}

type orderRow struct {
	ID                 string
	SourceQuoteID      string
	CustomerID         string
	Status             string
	RequestedBy        string
	PaymentID          string
	PaymentStatus      string
	ShippedAt          time.Time
	CancelledBy        string
	CancellationReason string
	Lines              []OrderLine
	Total              int
}

type stockRow struct {
	SKU              string
	OnHand           int
	Reserved         int
	ReorderThreshold int
}

type paymentRow struct {
	ID              string
	OrderID         string
	Amount          int
	Status          string
	ReviewedBy      string
	DecisionComment string
}

type shipmentRow struct {
	ID        string
	OrderID   string
	Status    string
	ShippedBy string
	ShippedAt time.Time
	Lines     []ShipmentLine
}

type returnRow struct {
	ID           string
	OrderID      string
	Status       string
	Reason       string
	ReviewNote   string
	RequestedBy  string
	ReviewedBy   string
	ProcessedBy  string
	Lines        []ReturnLine
	RefundID     string
	RefundStatus string
	RefundAmount int
	RequestedAt  time.Time
}

type refundRow struct {
	ID              string
	ReturnRequestID string
	OrderID         string
	Amount          int
	Status          string
	ProcessedBy     string
}

type pluginRow struct {
	Key     string
	Type    string
	Version string
	Enabled bool
	Config  map[string]string
}

// Database is the small persistence boundary used by this lesson. Its tables
// stay private so callers must use the Active Record operations instead of
// reaching into storage directly.
type Database struct {
	customers          map[string]customerRow
	quotes             map[string]quoteRow
	products           map[string]productRow
	orders             map[string]orderRow
	stocks             map[string]stockRow
	payments           map[string]paymentRow
	shipments          map[string]shipmentRow
	returns            map[string]returnRow
	refunds            map[string]refundRow
	plugins            map[string]pluginRow
	idempotency        map[string]string
	nextQuoteNumber    int
	nextOrderNumber    int
	nextPaymentNumber  int
	nextShipmentNumber int
	nextReturnNumber   int
	nextRefundNumber   int
}

func NewDatabase() *Database {
	return &Database{
		customers:   make(map[string]customerRow),
		quotes:      make(map[string]quoteRow),
		products:    make(map[string]productRow),
		orders:      make(map[string]orderRow),
		stocks:      make(map[string]stockRow),
		payments:    make(map[string]paymentRow),
		shipments:   make(map[string]shipmentRow),
		returns:     make(map[string]returnRow),
		refunds:     make(map[string]refundRow),
		plugins:     make(map[string]pluginRow),
		idempotency: make(map[string]string),
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

func (db *Database) nextPaymentID() string {
	db.nextPaymentNumber++
	return fmt.Sprintf("payment-%03d", db.nextPaymentNumber)
}

func (db *Database) nextShipmentID() string {
	db.nextShipmentNumber++
	return fmt.Sprintf("shipment-%03d", db.nextShipmentNumber)
}

func (db *Database) nextReturnID() string {
	db.nextReturnNumber++
	return fmt.Sprintf("return-%03d", db.nextReturnNumber)
}

func (db *Database) nextRefundID() string {
	db.nextRefundNumber++
	return fmt.Sprintf("refund-%03d", db.nextRefundNumber)
}
