package customers

import (
	"testing"

	"domain-driven-design-architecture/internal/domain/customer"
)

func TestReaderProjectsCustomerAggregate(t *testing.T) {
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
	if details.Tier != string(customer.CustomerTierPreferred) || !details.Active {
		t.Fatalf("unexpected details %+v", details)
	}
	if err := value.Deactivate(); err != nil {
		t.Fatal(err)
	}
	reader.Save(value)
	active := true
	if got := reader.ListCustomers(&active); len(got) != 0 {
		t.Fatalf("inactive customer appeared in active list: %+v", got)
	}
}
