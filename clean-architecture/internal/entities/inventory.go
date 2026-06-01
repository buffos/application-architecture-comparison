package entities

import "errors"

var ErrInsufficientInventory = errors.New("insufficient inventory")

type InventoryReservationItem struct {
	SKU      string
	Quantity int
}

type InventoryStockRecord struct {
	SKU      string
	Quantity int
}
