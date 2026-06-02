package quotes

import "testing"

func TestGetQuoteReturnsQuoteDetails(t *testing.T) {
	quotes := &stubQuoteRepository{
		saved: Quote{
			ID:         "quote-001",
			CustomerID: "customer-001",
			Status:     QuoteStatusDraft,
		},
	}

	service := NewService(quotes, stubCustomerDirectory{}, stubProductCatalog{})

	result, err := service.GetQuote(GetQuoteQuery{QuoteID: "quote-001"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.QuoteID != "quote-001" {
		t.Fatalf("expected quote-001, got %s", result.QuoteID)
	}

	if result.Status != QuoteStatusDraft {
		t.Fatalf("expected status %s, got %s", QuoteStatusDraft, result.Status)
	}

	if result.LineCount != 0 {
		t.Fatalf("expected line count 0, got %d", result.LineCount)
	}
}
