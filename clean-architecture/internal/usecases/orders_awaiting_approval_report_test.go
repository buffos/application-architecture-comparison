package usecases

import (
	"testing"

	"clean-architecture/internal/entities"
)

type stubOrdersAwaitingApprovalReportOutput struct {
	output OrdersAwaitingApprovalReportOutput
}

func (o *stubOrdersAwaitingApprovalReportOutput) Present(output OrdersAwaitingApprovalReportOutput) error {
	o.output = output
	return nil
}

func TestOrdersAwaitingApprovalReportInteractorBuildsQueue(t *testing.T) {
	output := &stubOrdersAwaitingApprovalReportOutput{}
	interactor := NewOrdersAwaitingApprovalReportInteractor(
		stubQuoteReportReader{
			byStatus: map[string][]entities.Quote{
				entities.QuoteStatusPendingApproval: {
					{
						ID:         "quote-001",
						CustomerID: "customer-001",
						Lines: []entities.QuoteLine{
							{LineTotal: 10000},
							{LineTotal: 2500},
						},
					},
				},
			},
		},
		output,
	)

	err := interactor.Execute(OrdersAwaitingApprovalReportInput{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output.output.Count != 1 {
		t.Fatalf("expected 1 queue item, got %d", output.output.Count)
	}

	if output.output.Items[0].QuoteID != "quote-001" {
		t.Fatalf("expected quote-001, got %s", output.output.Items[0].QuoteID)
	}

	if output.output.Items[0].TotalAmount != 12500 {
		t.Fatalf("expected total amount 12500, got %d", output.output.Items[0].TotalAmount)
	}
}
