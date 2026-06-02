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

	service := NewService(quotes, stubCustomerDirectory{}, stubProductCatalog{}, stubApprovalEvaluator{})

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

func TestListQuotesFiltersByStatus(t *testing.T) {
	quotes := &stubListQuoteRepository{
		quotes: map[string]Quote{
			"quote-001": {
				ID:         "quote-001",
				CustomerID: "customer-001",
				Status:     QuoteStatusDraft,
			},
			"quote-002": {
				ID:         "quote-002",
				CustomerID: "customer-002",
				Status:     QuoteStatusApproved,
				Lines: []QuoteLine{
					{ProductSKU: "sku-001", Quantity: 2, UnitPrice: 15000},
				},
			},
		},
	}

	service := NewService(quotes, stubCustomerDirectory{}, stubProductCatalog{}, stubApprovalEvaluator{})

	result, err := service.ListQuotes(ListQuotesQuery{Status: QuoteStatusApproved})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("expected one approved quote, got %d", len(result))
	}

	if result[0].QuoteID != "quote-002" || result[0].Status != QuoteStatusApproved {
		t.Fatalf("expected approved quote-002, got %+v", result[0])
	}
}

type stubListQuoteRepository struct {
	quotes map[string]Quote
}

func (r *stubListQuoteRepository) Save(quote Quote) error {
	if r.quotes == nil {
		r.quotes = make(map[string]Quote)
	}
	r.quotes[quote.ID] = quote
	return nil
}

func (r *stubListQuoteRepository) FindByID(id string) (Quote, error) {
	quote, ok := r.quotes[id]
	if !ok {
		return Quote{}, ErrQuoteNotFound
	}
	return quote, nil
}

func (r *stubListQuoteRepository) ListByStatus(status string) ([]Quote, error) {
	list := make([]Quote, 0, len(r.quotes))
	for _, quote := range r.quotes {
		if status == "" || quote.Status == status {
			list = append(list, quote)
		}
	}
	return list, nil
}
