package quotes

import (
	"testing"

	"microkernel-architecture/internal/kernel"
)

type stubRepository struct {
	saved Quote
}

func (r *stubRepository) FindByID(id string) (Quote, error) {
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

func TestCreateDraftQuote(t *testing.T) {
	repository := &stubRepository{}
	service := NewService(repository, stubCustomerDirectory{})

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
