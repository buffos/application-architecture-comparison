package memory

import (
	"sync"

	"onion-architecture/internal/domain"
)

type InventoryReservation struct {
	mu    sync.Mutex
	stock map[string]int
}

func NewInventoryReservation() *InventoryReservation {
	return &InventoryReservation{
		stock: make(map[string]int),
	}
}

func (r *InventoryReservation) Seed(sku string, quantity int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.stock[sku] = quantity
}

func (r *InventoryReservation) Reserve(items []domain.InventoryReservationItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, item := range items {
		if r.stock[item.ProductSKU] < item.Quantity {
			return domain.ErrInsufficientStock
		}
	}

	for _, item := range items {
		r.stock[item.ProductSKU] -= item.Quantity
	}

	return nil
}

func (r *InventoryReservation) Release(items []domain.InventoryReleaseItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, item := range items {
		r.stock[item.ProductSKU] += item.Quantity
	}

	return nil
}

func (r *InventoryReservation) Restock(items []domain.InventoryRestockItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, item := range items {
		r.stock[item.ProductSKU] += item.Quantity
	}

	return nil
}
