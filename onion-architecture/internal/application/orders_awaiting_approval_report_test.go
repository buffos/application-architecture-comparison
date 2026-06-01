package application

import (
	"testing"

	"onion-architecture/internal/domain"
)

func TestOrdersAwaitingApprovalReportServiceBuildsApprovalQueue(t *testing.T) {
	service := NewOrdersAwaitingApprovalReportService(stubQuoteListFinder{
		list: []domain.Quote{
			{
				ID:         "quote-001",
				CustomerID: "customer-001",
				Status:     domain.QuoteStatusPendingApproval,
				Lines: []domain.QuoteLine{
					{Quantity: 2, UnitPrice: 15000},
					{Quantity: 1, UnitPrice: 5000},
				},
			},
			{
				ID:         "quote-002",
				CustomerID: "customer-002",
				Status:     domain.QuoteStatusApproved,
				Lines: []domain.QuoteLine{
					{Quantity: 1, UnitPrice: 9999},
				},
			},
		},
	})

	report, err := service.Execute()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(report) != 1 {
		t.Fatalf("expected 1 queue row, got %d", len(report))
	}

	row := report[0]
	if row.QuoteID != "quote-001" {
		t.Fatalf("expected quote-001, got %s", row.QuoteID)
	}

	if row.LineCount != 2 {
		t.Fatalf("expected line count 2, got %d", row.LineCount)
	}

	if row.TotalAmount != 35000 {
		t.Fatalf("expected total amount 35000, got %d", row.TotalAmount)
	}
}
