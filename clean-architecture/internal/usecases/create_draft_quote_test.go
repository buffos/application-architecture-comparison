package usecases

import (
	"testing"

	"clean-architecture/internal/entities"
)

type stubQuoteGateway struct {
	saved entities.Quote
}

func (g *stubQuoteGateway) Save(quote entities.Quote) error {
	g.saved = quote
	return nil
}

type stubCustomerGateway struct {
	customer entities.Customer
	err      error
}

func (g stubCustomerGateway) FindByID(id string) (entities.Customer, error) {
	if g.err != nil {
		return entities.Customer{}, g.err
	}

	return g.customer, nil
}

type stubCreateDraftQuoteOutput struct {
	output CreateDraftQuoteOutput
}

func (o *stubCreateDraftQuoteOutput) Present(output CreateDraftQuoteOutput) error {
	o.output = output
	return nil
}

func TestCreateDraftQuoteInteractorCreatesDraftQuoteForActiveCustomer(t *testing.T) {
	quotes := &stubQuoteGateway{}
	customers := stubCustomerGateway{
		customer: entities.Customer{
			ID:     "customer-001",
			Active: true,
		},
	}
	output := &stubCreateDraftQuoteOutput{}

	interactor := NewCreateDraftQuoteInteractor(quotes, customers, output)

	err := interactor.Execute(CreateDraftQuoteInput{CustomerID: "customer-001"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if quotes.saved.CustomerID != "customer-001" {
		t.Fatalf("expected saved customer id customer-001, got %s", quotes.saved.CustomerID)
	}

	if quotes.saved.Status != entities.QuoteStatusDraft {
		t.Fatalf("expected saved status %s, got %s", entities.QuoteStatusDraft, quotes.saved.Status)
	}

	if output.output.QuoteID == "" {
		t.Fatal("expected presenter output to include quote id")
	}

	if output.output.Status != entities.QuoteStatusDraft {
		t.Fatalf("expected output status %s, got %s", entities.QuoteStatusDraft, output.output.Status)
	}
}

func TestCreateDraftQuoteInteractorRejectsInactiveCustomer(t *testing.T) {
	quotes := &stubQuoteGateway{}
	customers := stubCustomerGateway{
		customer: entities.Customer{
			ID:     "customer-002",
			Active: false,
		},
	}
	output := &stubCreateDraftQuoteOutput{}

	interactor := NewCreateDraftQuoteInteractor(quotes, customers, output)

	err := interactor.Execute(CreateDraftQuoteInput{CustomerID: "customer-002"})
	if err != entities.ErrCustomerInactive {
		t.Fatalf("expected %v, got %v", entities.ErrCustomerInactive, err)
	}

	if quotes.saved.ID != "" {
		t.Fatal("expected no quote to be saved for inactive customer")
	}
}
