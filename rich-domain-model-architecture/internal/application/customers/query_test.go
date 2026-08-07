package customers

import (
	"errors"
	"testing"

	"rich-domain-model-architecture/internal/domain/customer"
)

func TestCustomerQueryProjectsDetailsAndFiltersActiveState(t *testing.T) {
	value, err := customer.NewCustomer("customer-001", customer.CustomerTierPreferred, customer.PaymentTermsInvoice30)
	if err != nil {
		t.Fatal(err)
	}
	reader := NewInMemoryReader()
	reader.Save(value)
	details, err := reader.GetCustomer("customer-001")
	if err != nil {
		t.Fatal(err)
	}
	if details.Tier != string(customer.CustomerTierPreferred) || details.PaymentTerm != string(customer.PaymentTermsInvoice30) || !details.Active {
		t.Fatalf("details = %+v", details)
	}
	active := true
	rows := reader.ListCustomers(&active)
	if len(rows) != 1 || rows[0].CustomerID != "customer-001" {
		t.Fatalf("rows = %+v", rows)
	}
}

func TestCustomerQueryReturnsNotFound(t *testing.T) {
	if _, err := NewInMemoryReader().GetCustomer("missing"); !errors.Is(err, ErrCustomerNotFound) {
		t.Fatalf("missing query returned %v", err)
	}
}
