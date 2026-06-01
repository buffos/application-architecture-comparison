package memory

import (
	"sync"

	"clean-architecture/internal/entities"
)

type InventoryReservation struct {
	mu    sync.Mutex
	stock map[string]int
}

func NewInventoryReservation(stock map[string]int) *InventoryReservation {
	copyStock := make(map[string]int, len(stock))
	for sku, quantity := range stock {
		copyStock[sku] = quantity
	}

	return &InventoryReservation{
		stock: copyStock,
	}
}

func (r *InventoryReservation) Reserve(items []entities.InventoryReservationItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, item := range items {
		if r.stock[item.SKU] < item.Quantity {
			return entities.ErrInsufficientInventory
		}
	}

	for _, item := range items {
		r.stock[item.SKU] -= item.Quantity
	}

	return nil
}
