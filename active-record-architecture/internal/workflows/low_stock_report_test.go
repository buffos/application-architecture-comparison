package workflows

import (
	"testing"

	"active-record-architecture/internal/records"
)

func TestGetLowStockItemsReportCalculatesAvailableQuantity(t *testing.T) {
	db := records.NewDatabase()
	product := records.NewProduct(db, "sku-001", "Desk", "Standard", true, 15000)
	if err := product.Save(); err != nil {
		t.Fatalf("Product.Save() error = %v", err)
	}
	for _, stock := range []*records.StockRecord{
		records.NewStockRecord(db, "sku-002", 5, 3, 2),
		records.NewStockRecord(db, "sku-001", 10, 3, 7),
		records.NewStockRecord(db, "sku-003", 10, 2, 2),
	} {
		if err := stock.Save(); err != nil {
			t.Fatalf("StockRecord.Save() error = %v", err)
		}
	}

	items, err := records.GetLowStockItemsReport(db)
	if err != nil {
		t.Fatalf("GetLowStockItemsReport() error = %v", err)
	}
	if len(items) != 2 || items[0].SKU != "sku-001" || items[1].SKU != "sku-002" {
		t.Fatalf("low-stock items = %#v, want sorted sku-001 and sku-002", items)
	}
	if items[0].Available != 7 || items[0].ProductName != "Desk" || items[0].ReorderThreshold != 7 {
		t.Fatalf("threshold item = %#v", items[0])
	}
	if items[1].Available != 2 || items[1].ProductName != "" {
		t.Fatalf("missing-product item = %#v", items[1])
	}
}

func TestGetLowStockItemsReportExcludesHealthyStockAndRejectsMissingDatabase(t *testing.T) {
	db := records.NewDatabase()
	stock := records.NewStockRecord(db, "sku-001", 10, 0, 2)
	if err := stock.Save(); err != nil {
		t.Fatalf("StockRecord.Save() error = %v", err)
	}
	items, err := records.GetLowStockItemsReport(db)
	if err != nil {
		t.Fatalf("healthy report error = %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("healthy items = %#v, want empty", items)
	}
	if _, err := records.GetLowStockItemsReport(nil); err != records.ErrDatabaseRequired {
		t.Fatalf("missing database error = %v, want %v", err, records.ErrDatabaseRequired)
	}
}
