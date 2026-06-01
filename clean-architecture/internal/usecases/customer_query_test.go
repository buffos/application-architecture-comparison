package usecases

import (
	"testing"

	"clean-architecture/internal/entities"
)

type stubCustomerReader struct {
	customer entities.Customer
	err      error
}

func (g stubCustomerReader) FindByID(id string) (entities.Customer, error) {
	if g.err != nil {
		return entities.Customer{}, g.err
	}

	return g.customer, nil
}

type stubCustomerLister struct {
	customers  []entities.Customer
	err        error
	activeOnly bool
}

func (g *stubCustomerLister) List(activeOnly bool) ([]entities.Customer, error) {
	g.activeOnly = activeOnly
	if g.err != nil {
		return nil, g.err
	}

	return g.customers, nil
}

type stubGetCustomerOutput struct {
	output GetCustomerOutput
}

func (o *stubGetCustomerOutput) Present(output GetCustomerOutput) error {
	o.output = output
	return nil
}

type stubListCustomersOutput struct {
	output ListCustomersOutput
}

func (o *stubListCustomersOutput) Present(output ListCustomersOutput) error {
	o.output = output
	return nil
}

func TestGetCustomerInteractorLoadsCustomer(t *testing.T) {
	output := &stubGetCustomerOutput{}
	interactor := NewGetCustomerInteractor(stubCustomerReader{
		customer: entities.Customer{
			ID:     "customer-001",
			Active: true,
		},
	}, output)

	err := interactor.Execute(GetCustomerInput{CustomerID: "customer-001"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output.output.CustomerID != "customer-001" {
		t.Fatalf("expected customer id customer-001, got %s", output.output.CustomerID)
	}
}

func TestListCustomersInteractorFiltersActiveCustomers(t *testing.T) {
	customers := &stubCustomerLister{
		customers: []entities.Customer{
			{ID: "customer-001", Active: true},
			{ID: "customer-002", Active: true},
		},
	}
	output := &stubListCustomersOutput{}
	interactor := NewListCustomersInteractor(customers, output)

	err := interactor.Execute(ListCustomersInput{ActiveOnly: true})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !customers.activeOnly {
		t.Fatal("expected activeOnly filter to be true")
	}

	if output.output.Count != 2 {
		t.Fatalf("expected 2 customers, got %d", output.output.Count)
	}
}
