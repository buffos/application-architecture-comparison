package application

import (
	"hexagonal-architecture/internal/core/domain"
	"hexagonal-architecture/internal/core/ports"
)

type AdjustReorderThresholdUseCase struct {
	inventory ports.InventoryStockRepository
}

func NewAdjustReorderThresholdUseCase(inventory ports.InventoryStockRepository) AdjustReorderThresholdUseCase {
	return AdjustReorderThresholdUseCase{inventory: inventory}
}

func (uc AdjustReorderThresholdUseCase) Execute(sku string, threshold int) (domain.StockRecord, error) {
	record, err := uc.inventory.FindBySKU(sku)
	if err != nil {
		return domain.StockRecord{}, err
	}

	if err := record.SetReorderThreshold(threshold); err != nil {
		return domain.StockRecord{}, err
	}

	if err := uc.inventory.Save(record); err != nil {
		return domain.StockRecord{}, err
	}

	return record, nil
}
