package application

import "onion-architecture/internal/domain"

type InventoryStockReader interface {
	ListStock() ([]domain.InventoryStockRecord, error)
}

type LowStockItemsReportQuery struct {
	Threshold int
}

type LowStockItemRow struct {
	ProductSKU string
	Quantity   int
}

type LowStockItemsReportService struct {
	stock InventoryStockReader
}

func NewLowStockItemsReportService(stock InventoryStockReader) LowStockItemsReportService {
	return LowStockItemsReportService{
		stock: stock,
	}
}

func (s LowStockItemsReportService) Execute(query LowStockItemsReportQuery) ([]LowStockItemRow, error) {
	records, err := s.stock.ListStock()
	if err != nil {
		return nil, err
	}

	result := make([]LowStockItemRow, 0)
	for _, record := range records {
		if record.Quantity > query.Threshold {
			continue
		}

		result = append(result, LowStockItemRow{
			ProductSKU: record.ProductSKU,
			Quantity:   record.Quantity,
		})
	}

	return result, nil
}
