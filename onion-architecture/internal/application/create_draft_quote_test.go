package application

import (
	"testing"

	"onion-architecture/internal/domain"
)

type stubQuoteRepository struct {
	saved domain.Quote
}

func (r *stubQuoteRepository) Save(quote domain.Quote) error {
	r.saved = quote
	return nil
}

type stubCustomerRepository struct {
	customer domain.Customer
	err      error
}

func (r stubCustomerRepository) FindByID(id string) (domain.Customer, error) {
	if r.err != nil {
		return domain.Customer{}, r.err
	}

	return r.customer, nil
}

func TestCreateDraftQuoteServiceCreatesDraftQuoteForActiveCustomer(t *testing.T) {
	quotes := &stubQuoteRepository{}
	customers := stubCustomerRepository{
		customer: domain.Customer{
			ID:     "customer-001",
			Active: true,
		},
	}

	service := NewCreateDraftQuoteService(quotes, customers)

	result, err := service.Execute(CreateDraftQuoteCommand{CustomerID: "customer-001"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.CustomerID != "customer-001" {
		t.Fatalf("expected customer id customer-001, got %s", result.CustomerID)
	}

	if quotes.saved.Status != domain.QuoteStatusDraft {
		t.Fatalf("expected saved status %s, got %s", domain.QuoteStatusDraft, quotes.saved.Status)
	}
}

func TestCreateDraftQuoteServiceRejectsInactiveCustomer(t *testing.T) {
	quotes := &stubQuoteRepository{}
	customers := stubCustomerRepository{
		customer: domain.Customer{
			ID:     "customer-001",
			Active: false,
		},
	}

	service := NewCreateDraftQuoteService(quotes, customers)

	_, err := service.Execute(CreateDraftQuoteCommand{CustomerID: "customer-001"})
	if err != domain.ErrCustomerInactive {
		t.Fatalf("expected %v, got %v", domain.ErrCustomerInactive, err)
	}
}
