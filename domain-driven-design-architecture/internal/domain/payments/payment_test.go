package payments

import (
	"errors"
	"testing"
)

func TestPaymentAggregateCapturesOnce(t *testing.T) {
	amount, err := NewMoney(28500, "USD")
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

func TestPaymentAggregateCannotFailAfterCapture(t *testing.T) {
	amount, err := NewMoney(1000, "USD")
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
	if err := payment.Fail(); !errors.Is(err, ErrPaymentNotCapturable) {
		t.Fatalf("failure after capture returned %v", err)
	}
}
