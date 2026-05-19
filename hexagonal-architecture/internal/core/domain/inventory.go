package domain

import "errors"

var ErrInsufficientStock = errors.New("insufficient stock")
var ErrStockRecordNotFound = errors.New("stock record not found")
var ErrStockQuantityInvalid = errors.New("stock quantity must be positive")
var ErrReorderThresholdInvalid = errors.New("reorder threshold cannot be negative")

type StockRecord struct {
	SKU              string
	Available        int
	ReorderThreshold int
}

func (s *StockRecord) Receive(quantity int) error {
	if quantity <= 0 {
		return ErrStockQuantityInvalid
	}

	s.Available += quantity
	return nil
}

func (s *StockRecord) SetReorderThreshold(threshold int) error {
	if threshold < 0 {
		return ErrReorderThresholdInvalid
	}

	s.ReorderThreshold = threshold
	return nil
}

type ReservationLine struct {
	SKU      string
	Quantity int
}
