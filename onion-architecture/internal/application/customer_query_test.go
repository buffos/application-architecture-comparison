package application

import (
	"testing"

	"onion-architecture/internal/domain"
)

type stubCustomerLookup struct {
	customer domain.Customer
	list     []domain.Customer
	err      error
}

func (l stubCustomerLookup) FindByID(id string) (domain.Customer, error) {
	if l.err != nil {
		return domain.Customer{}, l.err
	}

	return l.customer, nil
}

func (l stubCustomerLookup) List(activeOnly bool) ([]domain.Customer, error) {
	if l.err != nil {
		return nil, l.err
	}

	result := make([]domain.Customer, 0)
	for _, customer := range l.list {
		if activeOnly && !customer.Active {
			continue
		}

		result = append(result, customer)
	}

	return result, nil
}

func TestGetCustomerServiceReturnsDetails(t *testing.T) {
	service := NewGetCustomerService(stubCustomerLookup{
		customer: domain.Customer{
			ID:     "customer-001",
			Active: true,
		},
	})

	result, err := service.Execute(GetCustomerQuery{CustomerID: "customer-001"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.CustomerID != "customer-001" {
		t.Fatalf("expected customer-001, got %s", result.CustomerID)
	}
}

func TestListCustomersServiceFiltersByActiveStatus(t *testing.T) {
	service := NewListCustomersService(stubCustomerLookup{
		list: []domain.Customer{
			{ID: "customer-001", Active: true},
			{ID: "customer-002", Active: false},
		},
	})

	result, err := service.Execute(ListCustomersQuery{ActiveOnly: true})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}

	if result[0].CustomerID != "customer-001" {
		t.Fatalf("expected customer-001, got %s", result[0].CustomerID)
	}
}
