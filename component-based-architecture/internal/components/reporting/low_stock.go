package reporting

import "component-based-architecture/internal/components/inventory"

type LowStockItem struct {
	ProductSKU string
	Available  int
}

type LowStockItemsReport struct {
	Threshold int
	Items     []LowStockItem
}

type LowStockComponent struct {
	inventory inventory.StockReader
}

func NewLowStockComponent(inventory inventory.StockReader) *LowStockComponent {
	return &LowStockComponent{inventory: inventory}
}

func (c *LowStockComponent) LowStockItemsReport(threshold int) LowStockItemsReport {
	report := LowStockItemsReport{Threshold: threshold, Items: make([]LowStockItem, 0)}
	for _, item := range c.inventory.ListStock() {
		if item.Available <= threshold {
			report.Items = append(report.Items, LowStockItem{ProductSKU: item.ProductSKU, Available: item.Available})
		}
	}
	return report
}
