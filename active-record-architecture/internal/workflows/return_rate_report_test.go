package workflows

import (
	"testing"

	"active-record-architecture/internal/records"
)

func TestGetReturnRateByCategoryReportAggregatesRows(t *testing.T) {
	db, first := readyOrder(t)
	first, err := records.FindOrder(db, first.ID)
	if err != nil {
		t.Fatalf("FindOrder() error = %v", err)
	}
	first.Status = records.OrderStatusShipped
	first.Lines[0].ProductCategory = "Standard"
	first.Lines[0].ShippedQuantity = 3
	first.Lines[0].ReturnedQuantity = 1
	if err := first.Save(); err != nil {
		t.Fatalf("Order.Save() error = %v", err)
	}

	second := secondReadyOrder(t, db)
	second, err = records.FindOrder(db, second.ID)
	if err != nil {
		t.Fatalf("second FindOrder() error = %v", err)
	}
	second.Status = records.OrderStatusShipped
	second.Lines[0].ProductCategory = "Clearance"
	second.Lines[0].ShippedQuantity = 2
	second.Lines[0].ReturnedQuantity = 0
	if err := second.Save(); err != nil {
		t.Fatalf("second Order.Save() error = %v", err)
	}

	report, err := records.GetReturnRateByCategoryReport(db)
	if err != nil {
		t.Fatalf("GetReturnRateByCategoryReport() error = %v", err)
	}
	if len(report.Rows) != 2 {
		t.Fatalf("report rows = %#v, want two categories", report.Rows)
	}
	if report.Rows[0].Category != "Clearance" || report.Rows[0].ShippedQuantity != 2 || report.Rows[0].ReturnedQuantity != 0 || report.Rows[0].ReturnRate != 0 {
		t.Fatalf("clearance row = %#v", report.Rows[0])
	}
	if report.Rows[1].Category != "Standard" || report.Rows[1].ShippedQuantity != 3 || report.Rows[1].ReturnedQuantity != 1 || report.Rows[1].ReturnRate != 1.0/3.0 {
		t.Fatalf("standard row = %#v", report.Rows[1])
	}
}

func TestGetReturnRateByCategoryReportHandlesEmptyAndUnknownCategories(t *testing.T) {
	empty, err := records.GetReturnRateByCategoryReport(records.NewDatabase())
	if err != nil {
		t.Fatalf("empty report error = %v", err)
	}
	if len(empty.Rows) != 0 {
		t.Fatalf("empty rows = %#v", empty.Rows)
	}

	db, order := readyOrder(t)
	order, err = records.FindOrder(db, order.ID)
	if err != nil {
		t.Fatalf("FindOrder() error = %v", err)
	}
	order.Status = records.OrderStatusShipped
	order.Lines[0].ProductCategory = ""
	order.Lines[0].ShippedQuantity = 1
	if err := order.Save(); err != nil {
		t.Fatalf("Order.Save() error = %v", err)
	}
	report, err := records.GetReturnRateByCategoryReport(db)
	if err != nil {
		t.Fatalf("unknown-category report error = %v", err)
	}
	if len(report.Rows) != 1 || report.Rows[0].Category != "Unknown" {
		t.Fatalf("unknown category rows = %#v", report.Rows)
	}
}
