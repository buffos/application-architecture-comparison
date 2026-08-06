package customers

import "testing"

func TestCustomerReaderLoadsAndFiltersCustomers(t *testing.T) {
	component := NewComponent()
	if err := component.Register(Customer{ID: "customer-001", Active: true}); err != nil {
		t.Fatal(err)
	}
	if err := component.Register(Customer{ID: "customer-002", Active: false}); err != nil {
		t.Fatal(err)
	}
	var reader Reader = component
	details, err := reader.GetCustomer(GetCustomerQuery{ID: "customer-001"})
	if err != nil || !details.Active {
		t.Fatalf("details=%+v err=%v", details, err)
	}
	active := reader.ListCustomers(ListCustomersQuery{ActiveOnly: true})
	if len(active) != 1 || active[0].ID != "customer-001" {
		t.Fatalf("unexpected active customers %+v", active)
	}
}
