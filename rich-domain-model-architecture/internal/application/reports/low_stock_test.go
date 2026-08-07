package reports

import (
	"testing"

	"rich-domain-model-architecture/internal/domain/inventory"
)

func TestLowStockReportUsesAggregateAvailability(t *testing.T) {
	first, err := inventory.NewStockRecord("sku-001", 10, 2)
	if err != nil {
		t.Fatal(err)
	}
	if err := first.Reserve(4); err != nil {
		t.Fatal(err)
	}
	second, err := inventory.NewStockRecord("sku-002", 20, 2)
	if err != nil {
		t.Fatal(err)
	}
	report := BuildLowStockItemsReport([]inventory.StockRecord{second, first}, 6)
	if len(report.Items) != 1 || report.Items[0].ProductSKU != "sku-001" || report.Items[0].Available != 6 {
		t.Fatalf("report = %+v", report)
	}
}

func TestLowStockReportReturnsEmptyForHighAvailability(t *testing.T) {
	record, err := inventory.NewStockRecord("sku-001", 10, 2)
	if err != nil {
		t.Fatal(err)
	}
	report := BuildLowStockItemsReport([]inventory.StockRecord{record}, 5)
	if len(report.Items) != 0 {
		t.Fatalf("report = %+v", report)
	}
}
