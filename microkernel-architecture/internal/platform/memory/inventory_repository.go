package memory

import (
	"sync"

	"microkernel-architecture/internal/plugins/inventory"
)

type InventoryRepository struct {
	mu    sync.RWMutex
	stock map[string]inventory.StockRecord
}

func NewInventoryRepository() *InventoryRepository {
	return &InventoryRepository{
		stock: make(map[string]inventory.StockRecord),
	}
}

func (r *InventoryRepository) Save(stock inventory.StockRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.stock[stock.ProductSKU] = stock
	return nil
}

func (r *InventoryRepository) FindBySKU(sku string) (inventory.StockRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.stock[sku], nil
}
