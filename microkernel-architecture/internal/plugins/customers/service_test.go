package customers

import "testing"

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
