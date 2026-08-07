package scripts

import (
	"testing"

	"transaction-script-architecture/internal/data"
)

func TestGetLowStockItemsReportCalculatesAvailableQuantity(t *testing.T) {
	store := data.NewStore()
	store.Products["sku-001"] = data.Product{SKU: "sku-001", Name: "Desk"}
	store.Products["sku-002"] = data.Product{SKU: "sku-002", Name: "Chair"}
	store.Stocks["sku-002"] = data.StockRecord{SKU: "sku-002", OnHand: 10, Reserved: 2, ReorderThreshold: 8}
	store.Stocks["sku-001"] = data.StockRecord{SKU: "sku-001", OnHand: 10, Reserved: 5, ReorderThreshold: 6}

	items, err := GetLowStockItemsReport(store)
	if err != nil {
		t.Fatalf("GetLowStockItemsReport() error = %v", err)
	}
	if len(items) != 2 || items[0].SKU != "sku-001" || items[1].SKU != "sku-002" {
		t.Fatalf("items = %#v, want sorted low-stock items", items)
	}
	if items[0].Available != 5 || items[0].ProductName != "Desk" {
		t.Fatalf("first item = %#v, want available 5 and Desk", items[0])
	}
}

func TestGetLowStockItemsReportExcludesHealthyStock(t *testing.T) {
	store := data.NewStore()
	store.Stocks["sku-001"] = data.StockRecord{SKU: "sku-001", OnHand: 10, Reserved: 1, ReorderThreshold: 5}

	items, err := GetLowStockItemsReport(store)
	if err != nil {
		t.Fatalf("GetLowStockItemsReport() error = %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("items = %#v, want empty", items)
	}
}
