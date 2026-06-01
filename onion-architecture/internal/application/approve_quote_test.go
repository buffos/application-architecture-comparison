package application

import (
	"testing"

	"onion-architecture/internal/domain"
)

func TestApproveQuoteServiceApprovesPendingQuote(t *testing.T) {
	quotes := &stubQuoteStore{
		quote: domain.Quote{
			ID:         "quote-001",
			CustomerID: "customer-001",
			Status:     domain.QuoteStatusPendingApproval,
			Lines: []domain.QuoteLine{
				{
					ProductSKU:      "sku-002",
					ProductName:     "Custom Desk",
					ProductCategory: "CustomBuild",
					Quantity:        1,
					UnitPrice:       45000,
				},
			},
		},
	}

	service := NewApproveQuoteService(quotes)

	result, err := service.Execute(ApproveQuoteCommand{QuoteID: "quote-001"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Status != domain.QuoteStatusApproved {
		t.Fatalf("expected status %s, got %s", domain.QuoteStatusApproved, result.Status)
	}
}

func TestApproveQuoteServiceRejectsAlreadyApprovedQuote(t *testing.T) {
	quotes := &stubQuoteStore{
		quote: domain.Quote{
			ID:         "quote-001",
			CustomerID: "customer-001",
			Status:     domain.QuoteStatusApproved,
			Lines: []domain.QuoteLine{
				{
					ProductSKU:      "sku-001",
					ProductName:     "Desk",
					ProductCategory: "Standard",
					Quantity:        1,
					UnitPrice:       15000,
				},
			},
		},
	}

	service := NewApproveQuoteService(quotes)

	_, err := service.Execute(ApproveQuoteCommand{QuoteID: "quote-001"})
	if err != domain.ErrQuoteNotApprovable {
		t.Fatalf("expected %v, got %v", domain.ErrQuoteNotApprovable, err)
	}
}
