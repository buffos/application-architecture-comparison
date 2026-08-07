package records

import "errors"

var (
	ErrStockSKURequired = errors.New("stock sku is required")
	ErrStockNotFound    = errors.New("stock record not found")
)

// StockRecord is an Active Record for inventory quantities.
type StockRecord struct {
	db *Database

	SKU              string
	OnHand           int
	Reserved         int
	ReorderThreshold int
}

// NewStockRecord creates a new, unsaved stock Active Record.
func NewStockRecord(db *Database, sku string, onHand int, reserved int, reorderThreshold int) *StockRecord {
	return &StockRecord{
		db:               db,
		SKU:              sku,
		OnHand:           onHand,
		Reserved:         reserved,
		ReorderThreshold: reorderThreshold,
	}
}

// FindStock loads an inventory Active Record from the stock table.
func FindStock(db *Database, sku string) (*StockRecord, error) {
	if db == nil {
		return nil, ErrDatabaseRequired
	}
	if sku == "" {
		return nil, ErrStockSKURequired
	}

	row, ok := db.stocks[sku]
	if !ok {
		return nil, ErrStockNotFound
	}
	return &StockRecord{
		db:               db,
		SKU:              row.SKU,
		OnHand:           row.OnHand,
		Reserved:         row.Reserved,
		ReorderThreshold: row.ReorderThreshold,
	}, nil
}

// Available returns the quantity not currently reserved.
func (stock *StockRecord) Available() int {
	if stock == nil {
		return 0
	}
	return stock.OnHand - stock.Reserved
}

// Reserve changes the row-level reserved quantity without saving it yet.
func (stock *StockRecord) Reserve(quantity int) error {
	if stock == nil || stock.db == nil {
		return ErrDatabaseRequired
	}
	if quantity <= 0 {
		return ErrQuantityInvalid
	}
	if stock.Available() < quantity {
		return ErrInsufficientStock
	}
	stock.Reserved += quantity
	return nil
}

// Consume removes shipped quantity from both on-hand and reserved stock.
func (stock *StockRecord) Consume(quantity int) error {
	if stock == nil || stock.db == nil {
		return ErrDatabaseRequired
	}
	if quantity <= 0 || stock.Reserved < quantity || stock.OnHand < quantity {
		return ErrInsufficientStock
	}
	stock.OnHand -= quantity
	stock.Reserved -= quantity
	return nil
}

// Save writes the current StockRecord Active Record to its table.
func (stock *StockRecord) Save() error {
	if stock == nil || stock.db == nil {
		return ErrDatabaseRequired
	}
	if stock.SKU == "" {
		return ErrStockSKURequired
	}
	stock.db.stocks[stock.SKU] = stockRow{
		SKU:              stock.SKU,
		OnHand:           stock.OnHand,
		Reserved:         stock.Reserved,
		ReorderThreshold: stock.ReorderThreshold,
	}
	return nil
}
