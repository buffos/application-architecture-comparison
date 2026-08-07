package payments

import (
	"errors"
	"testing"
)

func TestPaymentOwnsCaptureLifecycle(t *testing.T) {
	amount, err := NewMoney(30000, "USD")
	if err != nil {
		t.Fatal(err)
	}
	payment, err := NewPayment("payment-001", "order-001", amount)
	if err != nil {
		t.Fatal(err)
	}
	if err := payment.Capture(); err != nil {
		t.Fatal(err)
	}
	if payment.Status() != PaymentStatusCaptured {
		t.Fatalf("status = %s, want %s", payment.Status(), PaymentStatusCaptured)
	}
	if err := payment.Capture(); !errors.Is(err, ErrPaymentNotCapturable) {
		t.Fatalf("repeated capture returned %v", err)
	}
}

func TestPaymentRejectsInvalidIdentityAndTransitions(t *testing.T) {
	amount, err := NewMoney(1000, "USD")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := NewPayment("", "order-001", amount); !errors.Is(err, ErrPaymentIDRequired) {
		t.Fatalf("missing payment id returned %v", err)
	}
	if _, err := NewPayment("payment-001", "", amount); !errors.Is(err, ErrOrderIDRequired) {
		t.Fatalf("missing order id returned %v", err)
	}
	payment, err := NewPayment("payment-001", "order-001", amount)
	if err != nil {
		t.Fatal(err)
	}
	if err := payment.Fail(); err != nil {
		t.Fatal(err)
	}
	if err := payment.Capture(); !errors.Is(err, ErrPaymentNotCapturable) {
		t.Fatalf("capture after failure returned %v", err)
	}
}
