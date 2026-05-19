package domain

import "errors"

var ErrInsufficientStock = errors.New("insufficient stock")

type StockRecord struct {
	SKU              string
	Available        int
	ReorderThreshold int
}

type ReservationLine struct {
	SKU      string
	Quantity int
}
