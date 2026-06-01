package application

import (
	"testing"

	"onion-architecture/internal/domain"
)

func TestSubmitQuoteServiceSubmitsQuoteWithLines(t *testing.T) {
	quotes := &stubQuoteStore{
		quote: domain.Quote{
			ID:         "quote-001",
			CustomerID: "customer-001",
			Status:     domain.QuoteStatusDraft,
			Lines: []domain.QuoteLine{
				{
					ProductSKU:  "sku-001",
					ProductName: "Desk",
					Quantity:    2,
					UnitPrice:   15000,
				},
			},
		},
	}

	service := NewSubmitQuoteService(quotes)

	result, err := service.Execute(SubmitQuoteCommand{QuoteID: "quote-001"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Status != domain.QuoteStatusSubmitted {
		t.Fatalf("expected status %s, got %s", domain.QuoteStatusSubmitted, result.Status)
	}

	if quotes.saved.Status != domain.QuoteStatusSubmitted {
		t.Fatalf("expected saved status %s, got %s", domain.QuoteStatusSubmitted, quotes.saved.Status)
	}
}

func TestSubmitQuoteServiceRejectsEmptyQuote(t *testing.T) {
	quotes := &stubQuoteStore{
		quote: domain.Quote{
			ID:         "quote-001",
			CustomerID: "customer-001",
			Status:     domain.QuoteStatusDraft,
		},
	}

	service := NewSubmitQuoteService(quotes)

	_, err := service.Execute(SubmitQuoteCommand{QuoteID: "quote-001"})
	if err != domain.ErrQuoteCannotBeSubmittedWithoutLines {
		t.Fatalf("expected %v, got %v", domain.ErrQuoteCannotBeSubmittedWithoutLines, err)
	}
}
