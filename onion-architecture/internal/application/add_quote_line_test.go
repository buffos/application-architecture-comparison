package application

import (
	"testing"

	"onion-architecture/internal/domain"
)

type stubQuoteStore struct {
	quote domain.Quote
	err   error
	saved domain.Quote
}

func (s *stubQuoteStore) FindByID(id string) (domain.Quote, error) {
	if s.err != nil {
		return domain.Quote{}, s.err
	}

	return s.quote, nil
}

func (s *stubQuoteStore) Save(quote domain.Quote) error {
	s.saved = quote
	return nil
}

type stubProductLookup struct {
	product domain.Product
	err     error
}

func (s stubProductLookup) FindBySKU(sku string) (domain.Product, error) {
	if s.err != nil {
		return domain.Product{}, s.err
	}

	return s.product, nil
}

func TestAddQuoteLineServiceAddsLineToDraftQuote(t *testing.T) {
	quotes := &stubQuoteStore{
		quote: domain.Quote{
			ID:         "quote-001",
			CustomerID: "customer-001",
			Status:     domain.QuoteStatusDraft,
		},
	}

	products := stubProductLookup{
		product: domain.Product{
			SKU:      "sku-001",
			Name:     "Desk",
			Active:   true,
			UnitPrice: 15000,
		},
	}

	service := NewAddQuoteLineService(quotes, products)

	result, err := service.Execute(AddQuoteLineCommand{
		QuoteID:    "quote-001",
		ProductSKU: "sku-001",
		Quantity:   2,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.LineCount != 1 {
		t.Fatalf("expected line count 1, got %d", result.LineCount)
	}

	if result.TotalItems != 2 {
		t.Fatalf("expected total items 2, got %d", result.TotalItems)
	}

	if len(quotes.saved.Lines) != 1 {
		t.Fatalf("expected one saved line, got %d", len(quotes.saved.Lines))
	}
}

func TestAddQuoteLineServiceRejectsInactiveProduct(t *testing.T) {
	quotes := &stubQuoteStore{
		quote: domain.Quote{
			ID:         "quote-001",
			CustomerID: "customer-001",
			Status:     domain.QuoteStatusDraft,
		},
	}

	products := stubProductLookup{
		product: domain.Product{
			SKU:    "sku-002",
			Name:   "Legacy Desk",
			Active: false,
		},
	}

	service := NewAddQuoteLineService(quotes, products)

	_, err := service.Execute(AddQuoteLineCommand{
		QuoteID:    "quote-001",
		ProductSKU: "sku-002",
		Quantity:   1,
	})
	if err != domain.ErrProductInactive {
		t.Fatalf("expected %v, got %v", domain.ErrProductInactive, err)
	}
}
