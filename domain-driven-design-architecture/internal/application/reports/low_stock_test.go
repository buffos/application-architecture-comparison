package reports

import (
	"testing"

	"domain-driven-design-architecture/internal/domain/inventory"
)

func TestBuildLowStockItemsReport(t *testing.T) {
	low, err := inventory.NewStockRecord("sku-low", 5, 2)
	if err != nil {
		t.Fatal(err)
	}
	high, err := inventory.NewStockRecord("sku-high", 20, 2)
	if err != nil {
		t.Fatal(err)
	}
	report := BuildLowStockItemsReport([]inventory.StockRecord{high, low}, 5)
	if len(report.Items) != 1 || report.Items[0].ProductSKU != "sku-low" || report.Items[0].Available != 5 {
		t.Fatalf("unexpected report %+v", report)
	}
}
