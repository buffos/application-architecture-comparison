package quotes

import (
	"testing"

	"microkernel-architecture/internal/kernel"
)

type stubRepository struct {
	saved Quote
}

func (r *stubRepository) FindByID(id string) (Quote, error) {
	if r.saved.ID == id {
		return r.saved, nil
	}

	return Quote{}, ErrQuoteNotFound
}

func (r *stubRepository) Save(quote Quote) error {
	r.saved = quote
	return nil
}

type stubCustomerDirectory struct {
	err error
}

func (d stubCustomerDirectory) RequireActiveCustomer(id string) error {
	return d.err
}

type stubProductCatalog struct {
	product kernel.Product
	err     error
}

func (c stubProductCatalog) GetProductForQuote(sku string) (kernel.Product, error) {
	return c.product, c.err
}

func TestCreateDraftQuote(t *testing.T) {
	repository := &stubRepository{}
	service := NewService(repository, stubCustomerDirectory{}, stubProductCatalog{})

	result, err := service.CreateDraftQuote(kernel.CreateDraftQuoteCommand{
		CustomerID: "customer-001",
	})
	if err != nil {
		t.Fatalf("expected create draft quote to succeed, got %v", err)
	}

	if result.CustomerID != "customer-001" {
		t.Fatalf("expected customer id to be preserved, got %s", result.CustomerID)
	}

	if repository.saved.Status != QuoteStatusDraft {
		t.Fatalf("expected saved quote status %s, got %s", QuoteStatusDraft, repository.saved.Status)
	}
}

func TestGetQuote(t *testing.T) {
	repository := &stubRepository{
		saved: Quote{
			ID:         "quote-001",
			CustomerID: "customer-001",
			Status:     QuoteStatusDraft,
		},
	}
	service := NewService(repository, stubCustomerDirectory{}, stubProductCatalog{})

	result, err := service.GetQuote(kernel.GetQuoteQuery{
		QuoteID: "quote-001",
	})
	if err != nil {
		t.Fatalf("expected get quote to succeed, got %v", err)
	}

	if result.QuoteID != "quote-001" {
		t.Fatalf("expected quote id quote-001, got %s", result.QuoteID)
	}
}

func TestAddQuoteLine(t *testing.T) {
	repository := &stubRepository{
		saved: Quote{
			ID:         "quote-001",
			CustomerID: "customer-001",
			Status:     QuoteStatusDraft,
		},
	}
	service := NewService(repository, stubCustomerDirectory{}, stubProductCatalog{
		product: kernel.Product{
			SKU:       "sku-001",
			Name:      "Desk",
			UnitPrice: 15000,
		},
	})

	result, err := service.AddQuoteLine(kernel.AddQuoteLineCommand{
		QuoteID:    "quote-001",
		ProductSKU: "sku-001",
		Quantity:   2,
	})
	if err != nil {
		t.Fatalf("expected add quote line to succeed, got %v", err)
	}

	if result.LineCount != 1 {
		t.Fatalf("expected one line, got %d", result.LineCount)
	}

	if repository.saved.TotalQuantity() != 2 {
		t.Fatalf("expected total quantity 2, got %d", repository.saved.TotalQuantity())
	}
}
