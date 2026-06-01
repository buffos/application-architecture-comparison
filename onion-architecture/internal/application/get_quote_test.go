package application

import (
	"testing"

	"onion-architecture/internal/domain"
)

type stubQuoteFinder struct {
	quote domain.Quote
	err   error
}

func (f stubQuoteFinder) FindByID(id string) (domain.Quote, error) {
	if f.err != nil {
		return domain.Quote{}, f.err
	}

	return f.quote, nil
}

func (f stubQuoteFinder) ListByStatus(status string) ([]domain.Quote, error) {
	if f.err != nil {
		return nil, f.err
	}

	if f.quote.Status == status {
		return []domain.Quote{f.quote}, nil
	}

	return []domain.Quote{}, nil
}

func TestGetQuoteServiceReturnsQuoteDetails(t *testing.T) {
	service := NewGetQuoteService(stubQuoteFinder{
		quote: domain.Quote{
			ID:         "quote-001",
			CustomerID: "customer-001",
			Status:     domain.QuoteStatusDraft,
			Lines: []domain.QuoteLine{
				{
					ProductSKU:      "sku-001",
					ProductCategory: "Standard",
					Quantity:        2,
					UnitPrice:       15000,
				},
			},
		},
	})

	result, err := service.Execute(GetQuoteQuery{QuoteID: "quote-001"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.QuoteID != "quote-001" {
		t.Fatalf("expected quote id quote-001, got %s", result.QuoteID)
	}

	if result.Status != domain.QuoteStatusDraft {
		t.Fatalf("expected status %s, got %s", domain.QuoteStatusDraft, result.Status)
	}

	if result.LineCount != 1 {
		t.Fatalf("expected line count 1, got %d", result.LineCount)
	}
}

func TestGetQuoteServiceReturnsNotFoundError(t *testing.T) {
	service := NewGetQuoteService(stubQuoteFinder{
		err: domain.ErrQuoteNotFound,
	})

	_, err := service.Execute(GetQuoteQuery{QuoteID: "quote-404"})
	if err != domain.ErrQuoteNotFound {
		t.Fatalf("expected %v, got %v", domain.ErrQuoteNotFound, err)
	}
}
