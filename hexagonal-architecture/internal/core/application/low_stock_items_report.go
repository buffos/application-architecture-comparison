package application

import (
	"sort"

	"hexagonal-architecture/internal/core/ports"
)

type LowStockItemRow struct {
	SKU              string
	Available        int
	ReorderThreshold int
}

type GetLowStockItemsReportUseCase struct {
	inventory ports.InventoryStockReader
}

func NewGetLowStockItemsReportUseCase(inventory ports.InventoryStockReader) GetLowStockItemsReportUseCase {
	return GetLowStockItemsReportUseCase{inventory: inventory}
}

func (uc GetLowStockItemsReportUseCase) Execute() ([]LowStockItemRow, error) {
	records, err := uc.inventory.ListStock()
	if err != nil {
		return nil, err
	}

	rows := make([]LowStockItemRow, 0)
	for _, record := range records {
		if record.ReorderThreshold <= 0 {
			continue
		}
		if record.Available > record.ReorderThreshold {
			continue
		}

		rows = append(rows, LowStockItemRow{
			SKU:              record.SKU,
			Available:        record.Available,
			ReorderThreshold: record.ReorderThreshold,
		})
	}

	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Available == rows[j].Available {
			return rows[i].SKU < rows[j].SKU
		}
		return rows[i].Available < rows[j].Available
	})

	return rows, nil
}
