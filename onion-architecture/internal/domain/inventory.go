package domain

import "errors"

var ErrInsufficientStock = errors.New("insufficient stock")

type InventoryReservationItem struct {
	ProductSKU string
	Quantity   int
}

type InventoryReleaseItem struct {
	ProductSKU string
	Quantity   int
}

type InventoryRestockItem struct {
	ProductSKU string
	Quantity   int
}
