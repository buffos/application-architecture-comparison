package application

import (
	"hexagonal-architecture/internal/core/domain"
	"hexagonal-architecture/internal/core/ports"
)

type ReceiveStockUseCase struct {
	products  ports.ProductLookup
	inventory ports.InventoryStockRepository
}

func NewReceiveStockUseCase(products ports.ProductLookup, inventory ports.InventoryStockRepository) ReceiveStockUseCase {
	return ReceiveStockUseCase{
		products:  products,
		inventory: inventory,
	}
}

func (uc ReceiveStockUseCase) Execute(sku string, quantity int) (domain.StockRecord, error) {
	if _, err := uc.products.FindBySKU(sku); err != nil {
		return domain.StockRecord{}, err
	}

	record, err := uc.inventory.FindBySKU(sku)
	if err != nil {
		if err != domain.ErrStockRecordNotFound {
			return domain.StockRecord{}, err
		}
		record = domain.StockRecord{SKU: sku}
	}

	if err := record.Receive(quantity); err != nil {
		return domain.StockRecord{}, err
	}

	if err := uc.inventory.Save(record); err != nil {
		return domain.StockRecord{}, err
	}

	return record, nil
}
