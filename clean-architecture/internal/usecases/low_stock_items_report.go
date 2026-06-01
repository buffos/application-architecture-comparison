package usecases

import "clean-architecture/internal/entities"

type LowStockItemsReportInput struct {
	Threshold int
}

type LowStockItem struct {
	SKU      string
	Quantity int
}

type LowStockItemsReportOutput struct {
	Threshold int
	Count     int
	Items     []LowStockItem
}

type LowStockItemsReportInputBoundary interface {
	Execute(input LowStockItemsReportInput) error
}

type LowStockItemsReportOutputBoundary interface {
	Present(output LowStockItemsReportOutput) error
}

type InventoryStockReader interface {
	ListStock() ([]entities.InventoryStockRecord, error)
}

type LowStockItemsReportInteractor struct {
	inventory InventoryStockReader
	output    LowStockItemsReportOutputBoundary
}

func NewLowStockItemsReportInteractor(inventory InventoryStockReader, output LowStockItemsReportOutputBoundary) LowStockItemsReportInteractor {
	return LowStockItemsReportInteractor{
		inventory: inventory,
		output:    output,
	}
}

func (uc LowStockItemsReportInteractor) Execute(input LowStockItemsReportInput) error {
	records, err := uc.inventory.ListStock()
	if err != nil {
		return err
	}

	items := make([]LowStockItem, 0, len(records))
	for _, record := range records {
		if record.Quantity > input.Threshold {
			continue
		}

		items = append(items, LowStockItem{
			SKU:      record.SKU,
			Quantity: record.Quantity,
		})
	}

	return uc.output.Present(LowStockItemsReportOutput{
		Threshold: input.Threshold,
		Count:     len(items),
		Items:     items,
	})
}
