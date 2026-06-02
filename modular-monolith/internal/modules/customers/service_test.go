package customers

import "testing"

type stubCustomerRepository struct {
	customers map[string]Customer
}

func (r stubCustomerRepository) FindByID(id string) (Customer, error) {
	customer, ok := r.customers[id]
	if !ok {
		return Customer{}, ErrCustomerNotFound
	}
	return customer, nil
}

func (r stubCustomerRepository) List(activeOnly bool) ([]Customer, error) {
	list := make([]Customer, 0, len(r.customers))
	for _, customer := range r.customers {
		if activeOnly && !customer.Active {
			continue
		}
		list = append(list, customer)
	}
	return list, nil
}

func (r stubCustomerRepository) Save(customer Customer) error {
	return nil
}

func TestRequireActiveCustomerAcceptsActiveCustomer(t *testing.T) {
	service := NewService(stubCustomerRepository{
		customers: map[string]Customer{
			"customer-001": {ID: "customer-001", Active: true},
		},
	})

	err := service.RequireActiveCustomer("customer-001")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestRequireActiveCustomerRejectsInactiveCustomer(t *testing.T) {
	service := NewService(stubCustomerRepository{
		customers: map[string]Customer{
			"customer-001": {ID: "customer-001", Active: false},
		},
	})

	err := service.RequireActiveCustomer("customer-001")
	if err != ErrCustomerInactive {
		t.Fatalf("expected %v, got %v", ErrCustomerInactive, err)
	}
}

func TestGetCustomerReturnsStoredCustomer(t *testing.T) {
	service := NewService(stubCustomerRepository{
		customers: map[string]Customer{
			"customer-001": {ID: "customer-001", Active: true},
		},
	})

	customer, err := service.GetCustomer(GetCustomerQuery{CustomerID: "customer-001"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if customer.CustomerID != "customer-001" || !customer.Active {
		t.Fatalf("expected stored customer details, got %+v", customer)
	}
}

func TestListCustomersFiltersActiveCustomers(t *testing.T) {
	service := NewService(stubCustomerRepository{
		customers: map[string]Customer{
			"customer-001": {ID: "customer-001", Active: true},
			"customer-002": {ID: "customer-002", Active: false},
		},
	})

	result, err := service.ListCustomers(ListCustomersQuery{ActiveOnly: true})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result) != 1 || result[0].CustomerID != "customer-001" {
		t.Fatalf("expected one active customer, got %+v", result)
	}
}
