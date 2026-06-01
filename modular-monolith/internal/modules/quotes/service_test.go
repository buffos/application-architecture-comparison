package quotes

import (
	"testing"

	"modular-monolith/internal/modules/customers"
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

func TestCreateDraftQuoteCreatesDraftForActiveCustomer(t *testing.T) {
	quotes := &stubQuoteRepository{}
	service := NewService(quotes, stubCustomerDirectory{})

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
	})

	_, err := service.CreateDraftQuote(CreateDraftQuoteCommand{
		CustomerID: "customer-001",
	})
	if err != customers.ErrCustomerInactive {
		t.Fatalf("expected %v, got %v", customers.ErrCustomerInactive, err)
	}
}
