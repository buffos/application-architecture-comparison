package application

import (
	"hexagonal-architecture/internal/core/domain"
	"hexagonal-architecture/internal/core/ports"
)

type GetStockRecordUseCase struct {
	inventory ports.InventoryStockReader
}

func NewGetStockRecordUseCase(inventory ports.InventoryStockReader) GetStockRecordUseCase {
	return GetStockRecordUseCase{inventory: inventory}
}

func (uc GetStockRecordUseCase) Execute(sku string) (domain.StockRecord, error) {
	return uc.inventory.FindBySKU(sku)
}
