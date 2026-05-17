package domain

import "errors"

const StockShortageRejectOrder = "RejectOrder"

var ErrStockQuantityInvalid = errors.New("stock quantity must be positive")
var ErrStockRecordNotFound = errors.New("stock record not found")
var ErrInsufficientStock = errors.New("insufficient stock")

type StockRecord struct {
	SKU      string
	OnHand   int
	Reserved int
	Policy   string
}

func NewStockRecord(sku string, policy string) (StockRecord, error) {
	if sku == "" {
		return StockRecord{}, ErrProductSKURequired
	}

	if policy == "" {
		policy = StockShortageRejectOrder
	}

	return StockRecord{
		SKU:    sku,
		Policy: policy,
	}, nil
}

func (s StockRecord) Available() int {
	return s.OnHand - s.Reserved
}

func (s *StockRecord) Receive(quantity int) error {
	if quantity <= 0 {
		return ErrStockQuantityInvalid
	}

	s.OnHand += quantity
	return nil
}

func (s *StockRecord) Reserve(quantity int) error {
	if quantity <= 0 {
		return ErrStockQuantityInvalid
	}

	if s.Available() < quantity {
		return ErrInsufficientStock
	}

	s.Reserved += quantity
	return nil
}
