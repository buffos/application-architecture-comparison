package usecases

import (
	"testing"

	"clean-architecture/internal/entities"
)

type stubInventoryStockReader struct {
	records []entities.InventoryStockRecord
}

func (g stubInventoryStockReader) ListStock() ([]entities.InventoryStockRecord, error) {
	return g.records, nil
}

type stubLowStockItemsReportOutput struct {
	output LowStockItemsReportOutput
}

func (o *stubLowStockItemsReportOutput) Present(output LowStockItemsReportOutput) error {
	o.output = output
	return nil
}

func TestLowStockItemsReportInteractorFiltersByThreshold(t *testing.T) {
	output := &stubLowStockItemsReportOutput{}
	interactor := NewLowStockItemsReportInteractor(
		stubInventoryStockReader{
			records: []entities.InventoryStockRecord{
				{SKU: "CHAIR-001", Quantity: 1},
				{SKU: "DESK-001", Quantity: 5},
			},
		},
		output,
	)

	err := interactor.Execute(LowStockItemsReportInput{Threshold: 2})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output.output.Count != 1 {
		t.Fatalf("expected 1 low stock item, got %d", output.output.Count)
	}

	if output.output.Items[0].SKU != "CHAIR-001" {
		t.Fatalf("expected low stock item CHAIR-001, got %s", output.output.Items[0].SKU)
	}
}
