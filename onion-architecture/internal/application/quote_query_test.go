package application

import (
	"testing"

	"onion-architecture/internal/domain"
)

type stubQuoteListFinder struct {
	quote domain.Quote
	list  []domain.Quote
	err   error
}

func (f stubQuoteListFinder) FindByID(id string) (domain.Quote, error) {
	if f.err != nil {
		return domain.Quote{}, f.err
	}

	return f.quote, nil
}

func (f stubQuoteListFinder) ListByStatus(status string) ([]domain.Quote, error) {
	if f.err != nil {
		return nil, f.err
	}

	result := make([]domain.Quote, 0)
	for _, quote := range f.list {
		if quote.Status == status {
			result = append(result, quote)
		}
	}

	return result, nil
}

func TestListQuotesServiceFiltersByStatus(t *testing.T) {
	service := NewListQuotesService(stubQuoteListFinder{
		list: []domain.Quote{
			{ID: "quote-001", Status: domain.QuoteStatusDraft},
			{ID: "quote-002", Status: domain.QuoteStatusApproved},
		},
	})

	result, err := service.Execute(ListQuotesQuery{Status: domain.QuoteStatusApproved})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}

	if result[0].QuoteID != "quote-002" {
		t.Fatalf("expected quote-002, got %s", result[0].QuoteID)
	}
}
