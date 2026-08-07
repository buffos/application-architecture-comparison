package customer

import (
	"errors"
	"testing"
)

func TestCustomerOwnsLifecycleAndCommercialClassifications(t *testing.T) {
	customer, err := NewCustomer("customer-001", CustomerTierPreferred, PaymentTermsInvoice30)
	if err != nil {
		t.Fatal(err)
	}
	if !customer.Active() || customer.Tier() != CustomerTierPreferred || customer.PaymentTerms() != PaymentTermsInvoice30 {
		t.Fatalf("unexpected customer state: id=%s tier=%s terms=%s active=%t", customer.ID(), customer.Tier(), customer.PaymentTerms(), customer.Active())
	}
	if err := customer.EnsureCanCreateQuote(); err != nil {
		t.Fatal(err)
	}

	if err := customer.ChangeTier(CustomerTierEnterprise); err != nil {
		t.Fatal(err)
	}
	if err := customer.ChangePaymentTerms(PaymentTermsPrepaid); err != nil {
		t.Fatal(err)
	}
	if customer.Tier() != CustomerTierEnterprise || customer.PaymentTerms() != PaymentTermsPrepaid {
		t.Fatal("customer classification changes were not applied")
	}

	if err := customer.Deactivate(); err != nil {
		t.Fatal(err)
	}
	if err := customer.EnsureCanCreateQuote(); !errors.Is(err, ErrCustomerInactive) {
		t.Fatalf("inactive customer eligibility returned %v", err)
	}
	if err := customer.Deactivate(); !errors.Is(err, ErrCustomerAlreadyInactive) {
		t.Fatalf("repeated deactivation returned %v", err)
	}
	if err := customer.Activate(); err != nil {
		t.Fatal(err)
	}
	if err := customer.Activate(); !errors.Is(err, ErrCustomerAlreadyActive) {
		t.Fatalf("repeated activation returned %v", err)
	}
}

func TestCustomerRejectsInvalidIdentityAndClassifications(t *testing.T) {
	if _, err := NewCustomer("", CustomerTierStandard, PaymentTermsPrepaid); !errors.Is(err, ErrCustomerIDRequired) {
		t.Fatalf("empty identity returned %v", err)
	}
	if _, err := NewCustomer("customer-001", CustomerTier("Unknown"), PaymentTermsPrepaid); !errors.Is(err, ErrCustomerTierInvalid) {
		t.Fatalf("invalid tier returned %v", err)
	}
	if _, err := NewCustomer("customer-001", CustomerTierStandard, PaymentTerms("Unknown")); !errors.Is(err, ErrPaymentTermsInvalid) {
		t.Fatalf("invalid payment terms returned %v", err)
	}

	customer, err := NewCustomer("customer-001", CustomerTierStandard, PaymentTermsPrepaid)
	if err != nil {
		t.Fatal(err)
	}
	if err := customer.ChangeTier(CustomerTier("Unknown")); !errors.Is(err, ErrCustomerTierInvalid) {
		t.Fatalf("invalid tier change returned %v", err)
	}
	if err := customer.ChangePaymentTerms(PaymentTerms("Unknown")); !errors.Is(err, ErrPaymentTermsInvalid) {
		t.Fatalf("invalid payment terms change returned %v", err)
	}
}
