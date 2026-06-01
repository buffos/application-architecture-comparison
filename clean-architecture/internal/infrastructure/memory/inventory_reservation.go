package memory

import (
	"sort"
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

func (r *InventoryReservation) Release(items []entities.InventoryReservationItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, item := range items {
		r.stock[item.SKU] += item.Quantity
	}

	return nil
}

func (r *InventoryReservation) Restock(items []entities.InventoryReservationItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, item := range items {
		r.stock[item.SKU] += item.Quantity
	}

	return nil
}

func (r *InventoryReservation) ListStock() ([]entities.InventoryStockRecord, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	records := make([]entities.InventoryStockRecord, 0, len(r.stock))
	for sku, quantity := range r.stock {
		records = append(records, entities.InventoryStockRecord{
			SKU:      sku,
			Quantity: quantity,
		})
	}

	sort.Slice(records, func(i int, j int) bool {
		return records[i].SKU < records[j].SKU
	})

	return records, nil
}
