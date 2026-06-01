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

func TestGetQuoteServiceReturnsQuoteDetails(t *testing.T) {
	service := NewGetQuoteService(stubQuoteFinder{
		quote: domain.Quote{
			ID:         "quote-001",
			CustomerID: "customer-001",
			Status:     domain.QuoteStatusDraft,
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
