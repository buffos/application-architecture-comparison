package inventory

import "errors"

var (
	ErrStockNotFound                     = errors.New("stock record not found")
	ErrInsufficientStock                 = errors.New("insufficient stock")
	ErrReservationQuantityMustBePositive = errors.New("reservation quantity must be positive")
)

type StockRecord struct {
	ProductSKU string
	Available  int
}

type ReservationItem struct {
	ProductSKU string
	Quantity   int
}

type ReleaseItem struct {
	ProductSKU string
	Quantity   int
}
