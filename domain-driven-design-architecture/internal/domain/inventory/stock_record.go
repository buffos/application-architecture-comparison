package inventory

import "errors"

var (
	ErrProductSKURequired       = errors.New("product sku is required")
	ErrOnHandQuantityNegative   = errors.New("on-hand quantity cannot be negative")
	ErrReservedQuantityNegative = errors.New("reserved quantity cannot be negative")
	ErrReorderThresholdNegative = errors.New("reorder threshold cannot be negative")
	ErrQuantityMustBePositive   = errors.New("quantity must be positive")
	ErrInsufficientAvailable    = errors.New("insufficient available stock")
	ErrReleaseExceedsReserved   = errors.New("release exceeds reserved stock")
)

type ProductSKU string

// StockRecord is the aggregate root for the Inventory bounded context.
type StockRecord struct {
	sku              ProductSKU
	onHand           int
	reserved         int
	reorderThreshold int
}

func NewStockRecord(sku ProductSKU, onHand int, reorderThreshold int) (StockRecord, error) {
	if sku == "" {
		return StockRecord{}, ErrProductSKURequired
	}
	if onHand < 0 {
		return StockRecord{}, ErrOnHandQuantityNegative
	}
	if reorderThreshold < 0 {
		return StockRecord{}, ErrReorderThresholdNegative
	}
	return StockRecord{sku: sku, onHand: onHand, reorderThreshold: reorderThreshold}, nil
}

func (s StockRecord) SKU() ProductSKU       { return s.sku }
func (s StockRecord) OnHand() int           { return s.onHand }
func (s StockRecord) Reserved() int         { return s.reserved }
func (s StockRecord) Available() int        { return s.onHand - s.reserved }
func (s StockRecord) ReorderThreshold() int { return s.reorderThreshold }
func (s StockRecord) IsLowStock() bool      { return s.Available() <= s.reorderThreshold }

func (s *StockRecord) Receive(quantity int) error {
	if quantity <= 0 {
		return ErrQuantityMustBePositive
	}
	s.onHand += quantity
	return nil
}

func (s *StockRecord) Reserve(quantity int) error {
	if quantity <= 0 {
		return ErrQuantityMustBePositive
	}
	if quantity > s.Available() {
		return ErrInsufficientAvailable
	}
	s.reserved += quantity
	return nil
}

func (s *StockRecord) Release(quantity int) error {
	if quantity <= 0 {
		return ErrQuantityMustBePositive
	}
	if quantity > s.reserved {
		return ErrReleaseExceedsReserved
	}
	s.reserved -= quantity
	return nil
}
