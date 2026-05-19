package ports

import "hexagonal-architecture/internal/core/domain"

type InventoryStockReader interface {
	ListStock() ([]domain.StockRecord, error)
}
