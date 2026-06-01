package memory

import (
	"sort"
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

func (r *InventoryReservation) ListStock() ([]domain.InventoryStockRecord, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	result := make([]domain.InventoryStockRecord, 0, len(r.stock))
	for sku, quantity := range r.stock {
		result = append(result, domain.InventoryStockRecord{
			ProductSKU: sku,
			Quantity:   quantity,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ProductSKU < result[j].ProductSKU
	})

	return result, nil
}
