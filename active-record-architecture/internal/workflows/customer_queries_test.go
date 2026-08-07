package workflows

import (
	"testing"

	"active-record-architecture/internal/records"
)

func TestGetCustomerReturnsCustomerSnapshot(t *testing.T) {
	db := records.NewDatabase()
	customer := records.NewCustomer(db, "customer-001", true)
	if err := customer.Save(); err != nil {
		t.Fatalf("Customer.Save() error = %v", err)
	}

	got, err := records.GetCustomer(db, customer.ID)
	if err != nil {
		t.Fatalf("GetCustomer() error = %v", err)
	}
	got.Active = false
	saved, err := records.FindCustomer(db, customer.ID)
	if err != nil {
		t.Fatalf("FindCustomer() error = %v", err)
	}
	if !saved.Active {
		t.Fatalf("stored customer changed through query snapshot")
	}
}

func TestListCustomersFiltersAndSorts(t *testing.T) {
	db := records.NewDatabase()
	customers := []*records.Customer{
		records.NewCustomer(db, "customer-002", false),
		records.NewCustomer(db, "customer-001", true),
	}
	for _, customer := range customers {
		if err := customer.Save(); err != nil {
			t.Fatalf("Customer.Save() error = %v", err)
		}
	}

	active, err := records.ListCustomers(db, true)
	if err != nil {
		t.Fatalf("ListCustomers() active error = %v", err)
	}
	if len(active) != 1 || active[0].ID != "customer-001" {
		t.Fatalf("active customers = %#v, want customer-001", active)
	}

	all, err := records.ListCustomers(db, false)
	if err != nil {
		t.Fatalf("ListCustomers() all error = %v", err)
	}
	if len(all) != 2 || all[0].ID != "customer-001" || all[1].ID != "customer-002" {
		t.Fatalf("all customers = %#v, want sorted IDs", all)
	}
}

func TestGetCustomerRejectsMissingID(t *testing.T) {
	db := records.NewDatabase()
	if _, err := records.GetCustomer(db, "customer-404"); err != records.ErrCustomerNotFound {
		t.Fatalf("error = %v, want %v", err, records.ErrCustomerNotFound)
	}
}
