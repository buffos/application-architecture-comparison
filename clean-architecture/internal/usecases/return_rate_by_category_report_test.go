package usecases

import (
	"testing"

	"clean-architecture/internal/entities"
)

type stubReturnReportRequestReader struct {
	byStatus map[string][]entities.ReturnRequest
}

func (g stubReturnReportRequestReader) ListByStatus(status string) ([]entities.ReturnRequest, error) {
	return g.byStatus[status], nil
}

type stubReturnReportProductReader struct {
	bySKU map[string]entities.Product
}

func (g stubReturnReportProductReader) FindBySKU(sku string) (entities.Product, error) {
	return g.bySKU[sku], nil
}

type stubReturnRateByCategoryReportOutput struct {
	output ReturnRateByCategoryReportOutput
}

func (o *stubReturnRateByCategoryReportOutput) Present(output ReturnRateByCategoryReportOutput) error {
	o.output = output
	return nil
}

func TestReturnRateByCategoryReportInteractorBuildsProjection(t *testing.T) {
	output := &stubReturnRateByCategoryReportOutput{}
	interactor := NewReturnRateByCategoryReportInteractor(
		stubOrderReportReader{
			byStatus: map[string][]entities.Order{
				entities.OrderStatusShipped: {
					{
						ID: "order-001",
						Lines: []entities.OrderLine{
							{SKU: "CHAIR-001", Quantity: 2},
							{SKU: "DESK-001", Quantity: 1},
						},
					},
				},
			},
		},
		stubReturnReportRequestReader{
			byStatus: map[string][]entities.ReturnRequest{
				entities.ReturnRequestStatusRefunded: {
					{ID: "return-001", OrderID: "order-001"},
				},
			},
		},
		stubReturnReportProductReader{
			bySKU: map[string]entities.Product{
				"CHAIR-001": {SKU: "CHAIR-001", Category: "Seating"},
				"DESK-001":  {SKU: "DESK-001", Category: "Workspace"},
			},
		},
		output,
	)

	err := interactor.Execute(ReturnRateByCategoryReportInput{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(output.output.Categories) != 2 {
		t.Fatalf("expected 2 categories, got %d", len(output.output.Categories))
	}

	if output.output.Categories[0].Category != "Seating" {
		t.Fatalf("expected first category Seating, got %s", output.output.Categories[0].Category)
	}

	if output.output.Categories[0].ReturnRate != 1 {
		t.Fatalf("expected Seating return rate 1.00, got %f", output.output.Categories[0].ReturnRate)
	}
}
