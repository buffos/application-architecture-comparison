package ports

import "hexagonal-architecture/internal/core/domain"

type InventoryStockRepository interface {
	FindBySKU(sku string) (domain.StockRecord, error)
	Save(record domain.StockRecord) error
}
