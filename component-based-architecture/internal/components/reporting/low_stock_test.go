package reporting

import (
	"component-based-architecture/internal/components/inventory"
	"testing"
)

type stockReaderStub struct{}

func (stockReaderStub) ListStock() []inventory.StockSnapshot {
	return []inventory.StockSnapshot{{ProductSKU: "sku-001", Available: 2}, {ProductSKU: "sku-002", Available: 5}}
}

func TestLowStockItemsReportIncludesThresholdBoundary(t *testing.T) {
	report := NewLowStockComponent(stockReaderStub{}).LowStockItemsReport(2)
	if len(report.Items) != 1 || report.Items[0].ProductSKU != "sku-001" || report.Items[0].Available != 2 {
		t.Fatalf("unexpected report %+v", report)
	}
}
