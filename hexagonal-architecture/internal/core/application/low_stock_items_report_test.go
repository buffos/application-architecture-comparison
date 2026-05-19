package application

import (
	"testing"

	"hexagonal-architecture/internal/adapters/repository/memory"
)

func TestLowStockItemsReport(t *testing.T) {
	inventory := memory.NewInventoryReservationAdapter(map[string]int{
		"CHAIR-001": 2,
		"DESK-001":  6,
		"LAMP-001":  1,
		"RUG-001":   4,
	})
	inventory.SetReorderThreshold("CHAIR-001", 3)
	inventory.SetReorderThreshold("DESK-001", 5)
	inventory.SetReorderThreshold("LAMP-001", 1)

	reportUseCase := NewGetLowStockItemsReportUseCase(inventory)

	report, err := reportUseCase.Execute()
	if err != nil {
		t.Fatalf("expected report to succeed, got %v", err)
	}

	if len(report) != 2 {
		t.Fatalf("expected 2 low-stock rows, got %d", len(report))
	}

	if report[0].SKU != "LAMP-001" || report[0].Available != 1 || report[0].ReorderThreshold != 1 {
		t.Fatalf("unexpected first row: %+v", report[0])
	}

	if report[1].SKU != "CHAIR-001" || report[1].Available != 2 || report[1].ReorderThreshold != 3 {
		t.Fatalf("unexpected second row: %+v", report[1])
	}
}
