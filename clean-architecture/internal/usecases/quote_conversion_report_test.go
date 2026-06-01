package usecases

import (
	"testing"

	"clean-architecture/internal/entities"
)

type stubQuoteReportReader struct {
	byStatus map[string][]entities.Quote
}

func (g stubQuoteReportReader) ListByStatus(status string) ([]entities.Quote, error) {
	return g.byStatus[status], nil
}

type stubOrderReportReader struct {
	byStatus map[string][]entities.Order
}

func (g stubOrderReportReader) ListByStatus(status string) ([]entities.Order, error) {
	return g.byStatus[status], nil
}

type stubQuoteConversionReportOutput struct {
	output QuoteConversionReportOutput
}

func (o *stubQuoteConversionReportOutput) Present(output QuoteConversionReportOutput) error {
	o.output = output
	return nil
}

func TestQuoteConversionReportInteractorBuildsProjection(t *testing.T) {
	output := &stubQuoteConversionReportOutput{}
	interactor := NewQuoteConversionReportInteractor(
		stubQuoteReportReader{
			byStatus: map[string][]entities.Quote{
				entities.QuoteStatusDraft: {
					{ID: "quote-001"},
				},
				entities.QuoteStatusPendingApproval: {
					{ID: "quote-002"},
				},
				entities.QuoteStatusApproved: {
					{ID: "quote-003"},
					{ID: "quote-004"},
				},
			},
		},
		stubOrderReportReader{
			byStatus: map[string][]entities.Order{
				entities.OrderStatusPendingPayment: {
					{ID: "order-001"},
				},
				entities.OrderStatusPaid: {
					{ID: "order-002"},
				},
				entities.OrderStatusShipped: {
					{ID: "order-003"},
				},
				entities.OrderStatusPartiallyShipped: {},
				entities.OrderStatusCancelled:        {},
			},
		},
		output,
	)

	err := interactor.Execute(QuoteConversionReportInput{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output.output.TotalQuotes != 4 {
		t.Fatalf("expected 4 total quotes, got %d", output.output.TotalQuotes)
	}

	if output.output.ApprovedQuotes != 2 {
		t.Fatalf("expected 2 approved quotes, got %d", output.output.ApprovedQuotes)
	}

	if output.output.ConvertedQuotes != 3 {
		t.Fatalf("expected 3 converted quotes, got %d", output.output.ConvertedQuotes)
	}

	if output.output.ConversionRate != 0.75 {
		t.Fatalf("expected conversion rate 0.75, got %f", output.output.ConversionRate)
	}
}
