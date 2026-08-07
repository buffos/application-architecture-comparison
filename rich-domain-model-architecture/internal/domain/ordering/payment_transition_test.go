package ordering

import (
	"errors"
	"testing"
)

func TestOrderBecomesPaidOnlyFromPendingPayment(t *testing.T) {
	order, err := NewOrderFromQuote("order-001", approvedQuote(t))
	if err != nil {
		t.Fatal(err)
	}
	if err := order.MarkPaid(); err != nil {
		t.Fatal(err)
	}
	if order.Status() != OrderStatusPaid {
		t.Fatalf("status = %s, want %s", order.Status(), OrderStatusPaid)
	}
	if err := order.MarkPaid(); !errors.Is(err, ErrOrderNotAwaitingPayment) {
		t.Fatalf("repeated paid transition returned %v", err)
	}
}
