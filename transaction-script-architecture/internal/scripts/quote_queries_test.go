package scripts

import (
	"testing"

	"transaction-script-architecture/internal/data"
)

func TestGetQuoteReturnsDefensiveSnapshot(t *testing.T) {
	store := data.NewStore()
	store.Quotes["quote-001"] = data.Quote{
		ID:     "quote-001",
		Status: data.QuoteStatusDraft,
		Lines:  []data.QuoteLine{{SKU: "sku-001", Quantity: 1}},
	}

	got, err := GetQuote(store, "quote-001")
	if err != nil {
		t.Fatalf("GetQuote() error = %v", err)
	}
	got.Lines[0].Quantity = 99
	if store.Quotes["quote-001"].Lines[0].Quantity != 1 {
		t.Fatalf("store quantity = %d, want 1 after snapshot mutation", store.Quotes["quote-001"].Lines[0].Quantity)
	}
}

func TestListQuotesFiltersByStatusAndCustomer(t *testing.T) {
	store := data.NewStore()
	store.Quotes["quote-002"] = data.Quote{ID: "quote-002", CustomerID: "customer-002", Status: data.QuoteStatusApproved}
	store.Quotes["quote-001"] = data.Quote{ID: "quote-001", CustomerID: "customer-001", Status: data.QuoteStatusApproved}
	store.Quotes["quote-003"] = data.Quote{ID: "quote-003", CustomerID: "customer-001", Status: data.QuoteStatusDraft}

	got, err := ListQuotes(store, data.QuoteStatusApproved, "customer-001")
	if err != nil {
		t.Fatalf("ListQuotes() error = %v", err)
	}
	if len(got) != 1 || got[0].ID != "quote-001" {
		t.Fatalf("filtered quotes = %#v, want quote-001", got)
	}
}

func TestGetQuoteRejectsMissingQuote(t *testing.T) {
	store := data.NewStore()
	if _, err := GetQuote(store, "quote-404"); err != ErrQuoteNotFound {
		t.Fatalf("error = %v, want %v", err, ErrQuoteNotFound)
	}
}
