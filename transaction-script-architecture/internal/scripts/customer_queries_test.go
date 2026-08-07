package scripts

import (
	"testing"

	"transaction-script-architecture/internal/data"
)

func TestGetCustomerReturnsCustomer(t *testing.T) {
	store := data.NewStore()
	store.Customers["customer-001"] = data.Customer{ID: "customer-001", Active: true}

	got, err := GetCustomer(store, "customer-001")
	if err != nil {
		t.Fatalf("GetCustomer() error = %v", err)
	}
	if got.ID != "customer-001" || !got.Active {
		t.Fatalf("customer = %#v, want active customer-001", got)
	}
}

func TestListCustomersFiltersAndSorts(t *testing.T) {
	store := data.NewStore()
	store.Customers["customer-002"] = data.Customer{ID: "customer-002", Active: false}
	store.Customers["customer-001"] = data.Customer{ID: "customer-001", Active: true}

	active, err := ListCustomers(store, true)
	if err != nil {
		t.Fatalf("ListCustomers() active error = %v", err)
	}
	if len(active) != 1 || active[0].ID != "customer-001" {
		t.Fatalf("active customers = %#v, want customer-001", active)
	}

	all, err := ListCustomers(store, false)
	if err != nil {
		t.Fatalf("ListCustomers() all error = %v", err)
	}
	if len(all) != 2 || all[0].ID != "customer-001" || all[1].ID != "customer-002" {
		t.Fatalf("all customers = %#v, want sorted IDs", all)
	}
}

func TestGetCustomerRejectsMissingCustomer(t *testing.T) {
	store := data.NewStore()
	if _, err := GetCustomer(store, "customer-404"); err != ErrCustomerNotFound {
		t.Fatalf("error = %v, want %v", err, ErrCustomerNotFound)
	}
}
