package customer

import (
	"errors"
	"testing"
)

func TestCustomerAggregateOwnsLifecycle(t *testing.T) {
	customer, err := NewCustomer("customer-001", CustomerTierPreferred, PaymentTermsInvoice30)
	if err != nil {
		t.Fatal(err)
	}
	if !customer.Active() || customer.Tier() != CustomerTierPreferred || customer.PaymentTerms() != PaymentTermsInvoice30 {
		t.Fatalf("unexpected customer %+v", customer)
	}
	if err := customer.Deactivate(); err != nil {
		t.Fatal(err)
	}
	if customer.Active() {
		t.Fatal("customer should be inactive")
	}
	if err := customer.Deactivate(); !errors.Is(err, ErrCustomerAlreadyInactive) {
		t.Fatalf("repeated deactivation returned %v", err)
	}
	if err := customer.Activate(); err != nil {
		t.Fatal(err)
	}
	if !customer.Active() {
		t.Fatal("customer should be active")
	}
}

func TestCustomerAggregateRejectsInvalidIdentityAndClassifications(t *testing.T) {
	if _, err := NewCustomer("", CustomerTierStandard, PaymentTermsPrepaid); !errors.Is(err, ErrCustomerIDRequired) {
		t.Fatalf("empty identity returned %v", err)
	}
	if _, err := NewCustomer("customer-001", CustomerTier("Unknown"), PaymentTermsPrepaid); !errors.Is(err, ErrCustomerTierInvalid) {
		t.Fatalf("invalid tier returned %v", err)
	}
	if _, err := NewCustomer("customer-001", CustomerTierStandard, PaymentTerms("Unknown")); !errors.Is(err, ErrPaymentTermsInvalid) {
		t.Fatalf("invalid payment terms returned %v", err)
	}
}
