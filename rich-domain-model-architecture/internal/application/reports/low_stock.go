package reports

import (
	"sort"

	"rich-domain-model-architecture/internal/domain/inventory"
)

type LowStockItem struct {
	ProductSKU string
	Available  int
}

type LowStockItemsReport struct {
	Threshold int
	Items     []LowStockItem
}

func BuildLowStockItemsReport(records []inventory.StockRecord, threshold int) LowStockItemsReport {
	report := LowStockItemsReport{Threshold: threshold, Items: make([]LowStockItem, 0)}
	for _, record := range records {
		if record.Available() <= threshold {
			report.Items = append(report.Items, LowStockItem{ProductSKU: string(record.SKU()), Available: record.Available()})
		}
	}
	sort.Slice(report.Items, func(i, j int) bool { return report.Items[i].ProductSKU < report.Items[j].ProductSKU })
	return report
}
