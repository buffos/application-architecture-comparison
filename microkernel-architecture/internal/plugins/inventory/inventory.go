package inventory

import "errors"

var ErrInsufficientStock = errors.New("insufficient stock")

type StockRecord struct {
	ProductSKU string
	Available  int
}

type Repository interface {
	FindBySKU(sku string) (StockRecord, error)
	Save(stock StockRecord) error
}
