package customers

import (
	"testing"

	"microkernel-architecture/internal/kernel"
)

type stubRepository struct {
	customer Customer
	err      error
}

func (r stubRepository) FindByID(id string) (Customer, error) {
	return r.customer, r.err
}

func (r stubRepository) Save(customer Customer) error {
	return nil
}

func (r stubRepository) List(active *bool) ([]Customer, error) {
	if r.err != nil {
		return nil, r.err
	}

	if active != nil && r.customer.Active != *active {
		return []Customer{}, nil
	}

	return []Customer{r.customer}, nil
}

func TestRequireActiveCustomer(t *testing.T) {
	service := NewService(stubRepository{
		customer: Customer{
			ID:     "customer-001",
			Active: true,
		},
	})

	if err := service.RequireActiveCustomer("customer-001"); err != nil {
		t.Fatalf("expected active customer to pass, got %v", err)
	}
}

func TestRequireActiveCustomerRejectsInactive(t *testing.T) {
	service := NewService(stubRepository{
		customer: Customer{
			ID:     "customer-001",
			Active: false,
		},
	})

	if err := service.RequireActiveCustomer("customer-001"); err != ErrCustomerInactive {
		t.Fatalf("expected inactive error, got %v", err)
	}
}

func TestGetCustomer(t *testing.T) {
	service := NewService(stubRepository{
		customer: Customer{
			ID:     "customer-001",
			Active: true,
		},
	})

	customer, err := service.GetCustomer(kernel.GetCustomerQuery{CustomerID: "customer-001"})
	if err != nil {
		t.Fatalf("expected customer query to succeed, got %v", err)
	}

	if customer.CustomerID != "customer-001" || !customer.Active {
		t.Fatalf("unexpected customer details %+v", customer)
	}
}

func TestListCustomers(t *testing.T) {
	active := true
	service := NewService(stubRepository{
		customer: Customer{
			ID:     "customer-001",
			Active: true,
		},
	})

	customers, err := service.ListCustomers(kernel.ListCustomersQuery{
		Active: &active,
	})
	if err != nil {
		t.Fatalf("expected customer list to succeed, got %v", err)
	}

	if len(customers) != 1 || customers[0].CustomerID != "customer-001" {
		t.Fatalf("unexpected customer list %+v", customers)
	}
}
