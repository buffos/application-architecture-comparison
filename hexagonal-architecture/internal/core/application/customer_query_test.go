package application

import (
	"testing"

	"hexagonal-architecture/internal/adapters/repository/memory"
	"hexagonal-architecture/internal/core/domain"
)

func TestGetAndListCustomers(t *testing.T) {
	customerRepo := memory.NewCustomerRepository()
	getCustomer := NewGetCustomerUseCase(customerRepo)
	listCustomers := NewListCustomersUseCase(customerRepo)

	_ = customerRepo.Save(domain.Customer{ID: "customer-001", Active: true})
	_ = customerRepo.Save(domain.Customer{ID: "customer-002", Active: false})

	customer, err := getCustomer.Execute("customer-001")
	if err != nil {
		t.Fatalf("expected get customer to succeed, got %v", err)
	}

	if !customer.Active {
		t.Fatalf("expected active customer")
	}

	activeCustomers, err := listCustomers.Execute(true)
	if err != nil {
		t.Fatalf("expected list active customers to succeed, got %v", err)
	}

	if len(activeCustomers) != 1 || activeCustomers[0].ID != "customer-001" {
		t.Fatalf("expected one active customer customer-001, got %+v", activeCustomers)
	}

	allCustomers, err := listCustomers.Execute(false)
	if err != nil {
		t.Fatalf("expected list all customers to succeed, got %v", err)
	}

	if len(allCustomers) != 2 {
		t.Fatalf("expected two customers, got %+v", allCustomers)
	}
}
