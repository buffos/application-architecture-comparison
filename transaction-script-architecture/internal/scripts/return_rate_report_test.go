package scripts

import (
	"testing"

	"transaction-script-architecture/internal/data"
)

func TestGetReturnRateByCategoryReportAggregatesRows(t *testing.T) {
	store := data.NewStore()
	store.Orders["order-001"] = data.Order{
		ID: "order-001",
		Lines: []data.OrderLine{
			{ProductCategory: "Standard", ShippedQuantity: 4, ReturnedQuantity: 1},
			{ProductCategory: "Clearance", ShippedQuantity: 2, ReturnedQuantity: 0},
		},
	}
	store.Orders["order-002"] = data.Order{
		ID:    "order-002",
		Lines: []data.OrderLine{{ProductCategory: "Standard", ShippedQuantity: 1, ReturnedQuantity: 1}},
	}

	report, err := GetReturnRateByCategoryReport(store)
	if err != nil {
		t.Fatalf("GetReturnRateByCategoryReport() error = %v", err)
	}
	if len(report.Rows) != 2 || report.Rows[0].Category != "Clearance" || report.Rows[1].Category != "Standard" {
		t.Fatalf("rows = %#v, want sorted categories", report.Rows)
	}
	standard := report.Rows[1]
	if standard.ShippedQuantity != 5 || standard.ReturnedQuantity != 2 {
		t.Fatalf("standard row = %#v, want shipped 5 returned 2", standard)
	}
	if standard.ReturnRate < 0.399 || standard.ReturnRate > 0.401 {
		t.Fatalf("standard return rate = %f, want 0.4", standard.ReturnRate)
	}
	if report.Rows[0].ReturnRate != 0 {
		t.Fatalf("clearance return rate = %f, want 0", report.Rows[0].ReturnRate)
	}
}

func TestGetReturnRateByCategoryReportHandlesEmptyStore(t *testing.T) {
	report, err := GetReturnRateByCategoryReport(data.NewStore())
	if err != nil {
		t.Fatalf("GetReturnRateByCategoryReport() error = %v", err)
	}
	if len(report.Rows) != 0 {
		t.Fatalf("rows = %#v, want empty", report.Rows)
	}
}
