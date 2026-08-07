package quoting_test

import (
	"errors"
	"testing"

	"rich-domain-model-architecture/internal/domain/customer"
	"rich-domain-model-architecture/internal/domain/quoting"
)

func TestQuoteReferencesCustomerByIdentityOnly(t *testing.T) {
	customerAggregate, err := customer.NewCustomer("customer-001", customer.CustomerTierPreferred, customer.PaymentTermsInvoice30)
	if err != nil {
		t.Fatal(err)
	}
	if err := customerAggregate.EnsureCanCreateQuote(); err != nil {
		t.Fatal(err)
	}

	quote, err := quoting.NewQuote("quote-001", quoting.CustomerID(customerAggregate.ID()))
	if err != nil {
		t.Fatal(err)
	}
	if quote.CustomerID() != quoting.CustomerID(customerAggregate.ID()) {
		t.Fatalf("quote customer id = %s, want %s", quote.CustomerID(), customerAggregate.ID())
	}

	if err := customerAggregate.Deactivate(); err != nil {
		t.Fatal(err)
	}
	if err := customerAggregate.EnsureCanCreateQuote(); !errors.Is(err, customer.ErrCustomerInactive) {
		t.Fatalf("inactive customer eligibility returned %v", err)
	}
	if quote.CustomerID() != "customer-001" {
		t.Fatalf("existing quote reference changed to %s", quote.CustomerID())
	}
}
