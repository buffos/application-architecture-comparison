package quotes

import (
	"testing"

	"modular-monolith/internal/modules/customers"
	"modular-monolith/internal/modules/products"
)

type stubQuoteRepository struct {
	saved Quote
}

func (r *stubQuoteRepository) Save(quote Quote) error {
	r.saved = quote
	return nil
}

func (r *stubQuoteRepository) FindByID(id string) (Quote, error) {
	return r.saved, nil
}

type stubCustomerDirectory struct {
	err error
}

func (d stubCustomerDirectory) RequireActiveCustomer(id string) error {
	return d.err
}

type stubProductCatalog struct {
	product products.ProductForQuote
	err     error
}

func (c stubProductCatalog) GetProductForQuote(sku string) (products.ProductForQuote, error) {
	if c.err != nil {
		return products.ProductForQuote{}, c.err
	}

	return c.product, nil
}

func TestCreateDraftQuoteCreatesDraftForActiveCustomer(t *testing.T) {
	quotes := &stubQuoteRepository{}
	service := NewService(quotes, stubCustomerDirectory{}, stubProductCatalog{})

	result, err := service.CreateDraftQuote(CreateDraftQuoteCommand{
		CustomerID: "customer-001",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Status != QuoteStatusDraft {
		t.Fatalf("expected status %s, got %s", QuoteStatusDraft, result.Status)
	}

	if quotes.saved.CustomerID != "customer-001" {
		t.Fatalf("expected customer-001, got %s", quotes.saved.CustomerID)
	}
}

func TestCreateDraftQuoteRejectsInactiveCustomer(t *testing.T) {
	quotes := &stubQuoteRepository{}
	service := NewService(quotes, stubCustomerDirectory{
		err: customers.ErrCustomerInactive,
	}, stubProductCatalog{})

	_, err := service.CreateDraftQuote(CreateDraftQuoteCommand{
		CustomerID: "customer-001",
	})
	if err != customers.ErrCustomerInactive {
		t.Fatalf("expected %v, got %v", customers.ErrCustomerInactive, err)
	}
}

func TestAddQuoteLineAddsLineToExistingQuote(t *testing.T) {
	quotes := &stubQuoteRepository{
		saved: Quote{
			ID:         "quote-001",
			CustomerID: "customer-001",
			Status:     QuoteStatusDraft,
		},
	}
	service := NewService(quotes, stubCustomerDirectory{}, stubProductCatalog{
		product: products.ProductForQuote{
			SKU:       "sku-001",
			Name:      "Desk",
			Category:  "Standard",
			UnitPrice: 15000,
		},
	})

	result, err := service.AddQuoteLine(AddQuoteLineCommand{
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
}
