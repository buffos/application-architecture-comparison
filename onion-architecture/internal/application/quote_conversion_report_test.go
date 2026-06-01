package application

import (
	"testing"

	"onion-architecture/internal/domain"
)

func TestQuoteConversionReportServiceComputesMetrics(t *testing.T) {
	quotes := stubQuoteListFinder{
		list: []domain.Quote{
			{ID: "quote-001", Status: domain.QuoteStatusApproved},
			{ID: "quote-002", Status: domain.QuoteStatusPendingApproval},
			{ID: "quote-003", Status: domain.QuoteStatusDraft},
		},
	}

	orders := stubOrderFinder{
		list: []domain.Order{
			{ID: "order-001", Status: domain.OrderStatusPendingPayment},
			{ID: "order-002", Status: domain.OrderStatusShipped},
		},
	}

	service := NewQuoteConversionReportService(quotes, orders)

	report, err := service.Execute()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if report.TotalQuotes != 3 {
		t.Fatalf("expected total quotes 3, got %d", report.TotalQuotes)
	}

	if report.ApprovedQuotes != 2 {
		t.Fatalf("expected approved quotes 2, got %d", report.ApprovedQuotes)
	}

	if report.ConvertedQuotes != 2 {
		t.Fatalf("expected converted quotes 2, got %d", report.ConvertedQuotes)
	}
}
