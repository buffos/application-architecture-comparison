package quotes

import (
	"testing"

	"component-based-architecture/internal/components/customers"
)

func TestCreateDraftQuoteCreatesDraftForActiveCustomer(t *testing.T) {
	customerComponent := customers.NewComponent()
	if err := customerComponent.Register(customers.Customer{ID: "customer-001", Active: true}); err != nil {
		t.Fatalf("register customer: %v", err)
	}
	quoteComponent := NewComponent(customerComponent)

	result, err := quoteComponent.CreateDraftQuote(CreateDraftQuoteCommand{CustomerID: "customer-001"})
	if err != nil {
		t.Fatalf("create draft quote: %v", err)
	}

	if result.QuoteID != "quote-001" {
		t.Fatalf("expected quote-001, got %s", result.QuoteID)
	}
	if result.Status != QuoteStatusDraft {
		t.Fatalf("expected %s, got %s", QuoteStatusDraft, result.Status)
	}
}

func TestCreateDraftQuoteRejectsInactiveCustomer(t *testing.T) {
	customerComponent := customers.NewComponent()
	if err := customerComponent.Register(customers.Customer{ID: "customer-001", Active: false}); err != nil {
		t.Fatalf("register customer: %v", err)
	}
	quoteComponent := NewComponent(customerComponent)

	_, err := quoteComponent.CreateDraftQuote(CreateDraftQuoteCommand{CustomerID: "customer-001"})
	if err != customers.ErrCustomerInactive {
		t.Fatalf("expected %v, got %v", customers.ErrCustomerInactive, err)
	}
}
