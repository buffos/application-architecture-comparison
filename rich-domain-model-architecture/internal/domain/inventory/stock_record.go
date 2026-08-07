package inventory

import "errors"

var (
	ErrProductSKURequired       = errors.New("product sku is required")
	ErrOnHandQuantityNegative   = errors.New("on-hand quantity cannot be negative")
	ErrReorderThresholdNegative = errors.New("reorder threshold cannot be negative")
	ErrQuantityMustBePositive   = errors.New("quantity must be positive")
	ErrInsufficientAvailable    = errors.New("insufficient available stock")
	ErrReleaseExceedsReserved   = errors.New("release exceeds reserved stock")
)

type ProductSKU string

// StockRecord is the aggregate root for one SKU's inventory quantities.
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

func (record StockRecord) SKU() ProductSKU { return record.sku }
func (record StockRecord) OnHand() int { return record.onHand }
func (record StockRecord) Reserved() int { return record.reserved }
func (record StockRecord) Available() int { return record.onHand - record.reserved }
func (record StockRecord) ReorderThreshold() int { return record.reorderThreshold }
func (record StockRecord) IsLowStock() bool { return record.Available() <= record.reorderThreshold }

func (record *StockRecord) Receive(quantity int) error {
	if quantity <= 0 {
		return ErrQuantityMustBePositive
	}
	record.onHand += quantity
	return nil
}

func (record *StockRecord) Reserve(quantity int) error {
	if quantity <= 0 {
		return ErrQuantityMustBePositive
	}
	if quantity > record.Available() {
		return ErrInsufficientAvailable
	}
	record.reserved += quantity
	return nil
}

func (record *StockRecord) Release(quantity int) error {
	if quantity <= 0 {
		return ErrQuantityMustBePositive
	}
	if quantity > record.reserved {
		return ErrReleaseExceedsReserved
	}
	record.reserved -= quantity
	return nil
}
