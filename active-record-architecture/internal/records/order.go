package records

import "errors"

const (
	OrderStatusPendingReservation = "PendingReservation"
	PaymentStatusNotRequired      = "NotRequired"
)

var (
	ErrOrderIDRequired = errors.New("order id is required")
	ErrOrderNotFound   = errors.New("order not found")
)

// OrderLine is a committed product snapshot embedded in an Order Active
// Record.
type OrderLine struct {
	ID                  string
	SKU                 string
	ProductNameSnapshot string
	ProductCategory     string
	OrderedQuantity     int
	ReservedQuantity    int
	ShippedQuantity     int
	ReturnedQuantity    int
	UnitPrice           int
	DiscountAmount      int
	ReturnWindowDays    int
	LineTotal           int
}

// Order is an Active Record for a committed commercial transaction.
type Order struct {
	db *Database

	ID            string
	SourceQuoteID string
	CustomerID    string
	Status        string
	RequestedBy   string
	Lines         []OrderLine
	Total         int
	PaymentStatus string
}

// FindOrder loads an Order Active Record from the order table.
func FindOrder(db *Database, id string) (*Order, error) {
	if db == nil {
		return nil, ErrDatabaseRequired
	}
	if id == "" {
		return nil, ErrOrderIDRequired
	}

	row, ok := db.orders[id]
	if !ok {
		return nil, ErrOrderNotFound
	}

	return &Order{
		db:            db,
		ID:            row.ID,
		SourceQuoteID: row.SourceQuoteID,
		CustomerID:    row.CustomerID,
		Status:        row.Status,
		RequestedBy:   row.RequestedBy,
		PaymentStatus: row.PaymentStatus,
		Lines:         cloneOrderLines(row.Lines),
		Total:         row.Total,
	}, nil
}

// Save writes the current Order Active Record to its table.
func (order *Order) Save() error {
	if order == nil || order.db == nil {
		return ErrDatabaseRequired
	}
	if order.ID == "" {
		return ErrOrderIDRequired
	}

	order.db.orders[order.ID] = orderRow{
		ID:            order.ID,
		SourceQuoteID: order.SourceQuoteID,
		CustomerID:    order.CustomerID,
		Status:        order.Status,
		RequestedBy:   order.RequestedBy,
		PaymentStatus: order.PaymentStatus,
		Lines:         cloneOrderLines(order.Lines),
		Total:         order.Total,
	}
	return nil
}

func cloneOrderLines(lines []OrderLine) []OrderLine {
	if lines == nil {
		return nil
	}
	clone := make([]OrderLine, len(lines))
	copy(clone, lines)
	return clone
}
