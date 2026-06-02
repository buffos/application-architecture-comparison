package memory

import (
	"sync"

	"modular-monolith/internal/modules/inventory"
)

type InventoryRepository struct {
	mu      sync.Mutex
	records map[string]inventory.StockRecord
}

func NewInventoryRepository() *InventoryRepository {
	return &InventoryRepository{
		records: make(map[string]inventory.StockRecord),
	}
}

func (r *InventoryRepository) Save(record inventory.StockRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.records[record.ProductSKU] = record
	return nil
}

func (r *InventoryRepository) List() ([]inventory.StockRecord, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	list := make([]inventory.StockRecord, 0, len(r.records))
	for _, record := range r.records {
		list = append(list, record)
	}

	return list, nil
}

func (r *InventoryRepository) Reserve(items []inventory.ReservationItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, item := range items {
		record, ok := r.records[item.ProductSKU]
		if !ok {
			return inventory.ErrStockNotFound
		}

		if record.Available < item.Quantity {
			return inventory.ErrInsufficientStock
		}
	}

	for _, item := range items {
		record := r.records[item.ProductSKU]
		record.Available -= item.Quantity
		r.records[item.ProductSKU] = record
	}

	return nil
}

func (r *InventoryRepository) Release(items []inventory.ReleaseItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, item := range items {
		record, ok := r.records[item.ProductSKU]
		if !ok {
			return inventory.ErrStockNotFound
		}

		record.Available += item.Quantity
		r.records[item.ProductSKU] = record
	}

	return nil
}

func (r *InventoryRepository) Restock(items []inventory.RestockItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, item := range items {
		record, ok := r.records[item.ProductSKU]
		if !ok {
			return inventory.ErrStockNotFound
		}

		record.Available += item.Quantity
		r.records[item.ProductSKU] = record
	}

	return nil
}
