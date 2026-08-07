package scripts

import (
	"testing"

	"transaction-script-architecture/internal/data"
)

func TestAddQuoteLineAppendsAndPersistsLine(t *testing.T) {
	store := data.NewStore()
	store.Quotes["quote-001"] = data.Quote{
		ID:         "quote-001",
		CustomerID: "customer-001",
		Status:     data.QuoteStatusDraft,
	}
	store.Products["sku-001"] = data.Product{
		SKU:       "sku-001",
		Name:      "Desk",
		Category:  "Standard",
		Active:    true,
		UnitPrice: 15000,
	}

	got, err := AddQuoteLine(store, "quote-001", "sku-001", 2)
	if err != nil {
		t.Fatalf("AddQuoteLine() error = %v", err)
	}

	if len(got.Lines) != 1 {
		t.Fatalf("line count = %d, want 1", len(got.Lines))
	}

	line := got.Lines[0]
	if line.SKU != "sku-001" {
		t.Fatalf("line SKU = %q, want %q", line.SKU, "sku-001")
	}
	if line.ProductCategory != "Standard" {
		t.Fatalf("line product category = %q, want %q", line.ProductCategory, "Standard")
	}
	if line.ProductNameSnapshot != "Desk" {
		t.Fatalf("line product snapshot = %q, want %q", line.ProductNameSnapshot, "Desk")
	}
	if line.Quantity != 2 {
		t.Fatalf("line quantity = %d, want 2", line.Quantity)
	}
	if line.LineTotal != 30000 {
		t.Fatalf("line total = %d, want 30000", line.LineTotal)
	}

	saved := store.Quotes["quote-001"]
	if len(saved.Lines) != 1 || saved.Lines[0] != line {
		t.Fatalf("saved quote lines = %#v, want %#v", saved.Lines, got.Lines)
	}
}

func TestAddQuoteLineRejectsInvalidQuantity(t *testing.T) {
	store := data.NewStore()

	_, err := AddQuoteLine(store, "quote-001", "sku-001", 0)
	if err != ErrQuantityInvalid {
		t.Fatalf("error = %v, want %v", err, ErrQuantityInvalid)
	}
}

func TestAddQuoteLineRejectsUnknownQuote(t *testing.T) {
	store := data.NewStore()

	_, err := AddQuoteLine(store, "quote-404", "sku-001", 1)
	if err != ErrQuoteNotFound {
		t.Fatalf("error = %v, want %v", err, ErrQuoteNotFound)
	}
}

func TestAddQuoteLineRejectsUnknownProduct(t *testing.T) {
	store := data.NewStore()
	store.Quotes["quote-001"] = data.Quote{
		ID:     "quote-001",
		Status: data.QuoteStatusDraft,
	}

	_, err := AddQuoteLine(store, "quote-001", "sku-404", 1)
	if err != ErrProductNotFound {
		t.Fatalf("error = %v, want %v", err, ErrProductNotFound)
	}
}

func TestAddQuoteLineRejectsInactiveProduct(t *testing.T) {
	store := data.NewStore()
	store.Quotes["quote-001"] = data.Quote{
		ID:     "quote-001",
		Status: data.QuoteStatusDraft,
	}
	store.Products["sku-001"] = data.Product{
		SKU:    "sku-001",
		Active: false,
	}

	_, err := AddQuoteLine(store, "quote-001", "sku-001", 1)
	if err != ErrProductInactive {
		t.Fatalf("error = %v, want %v", err, ErrProductInactive)
	}
	if len(store.Quotes["quote-001"].Lines) != 0 {
		t.Fatalf("line count = %d, want 0", len(store.Quotes["quote-001"].Lines))
	}
}

func TestAddQuoteLineRejectsNonDraftQuote(t *testing.T) {
	store := data.NewStore()
	store.Quotes["quote-001"] = data.Quote{
		ID:     "quote-001",
		Status: "Submitted",
	}

	_, err := AddQuoteLine(store, "quote-001", "sku-001", 1)
	if err != ErrQuoteNotEditable {
		t.Fatalf("error = %v, want %v", err, ErrQuoteNotEditable)
	}
}
