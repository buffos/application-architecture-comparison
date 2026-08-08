package reporting

import (
	"testing"

	"rules-engine-architecture/internal/engine"
)

func TestBuildLowStockReportUsesStrictThreshold(t *testing.T) {
	rows := BuildLowStockReport([]engine.ProductFact{
		{ID: "PRD-003", Category: "Standard", AvailableQuantity: 0},
		{ID: "PRD-001", Category: "Standard", AvailableQuantity: 2},
		{ID: "PRD-002", Category: "CustomBuild", AvailableQuantity: 3},
	}, 3)

	if len(rows) != 2 || rows[0].ProductID != "PRD-001" || rows[1].ProductID != "PRD-003" {
		t.Fatalf("expected sorted low-stock rows below threshold, got %+v", rows)
	}
	if rows[0].Threshold != 3 {
		t.Fatalf("expected threshold in row, got %+v", rows[0])
	}
}

func TestBuildLowStockReportIgnoresNonPositiveThreshold(t *testing.T) {
	rows := BuildLowStockReport([]engine.ProductFact{{ID: "PRD-001", AvailableQuantity: 0}}, 0)

	if len(rows) != 0 {
		t.Fatalf("expected no rows for invalid threshold, got %+v", rows)
	}
}
