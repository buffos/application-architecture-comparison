package domain

import "errors"

var ErrInsufficientStock = errors.New("insufficient stock")

type ReservationLine struct {
	SKU      string
	Quantity int
}
