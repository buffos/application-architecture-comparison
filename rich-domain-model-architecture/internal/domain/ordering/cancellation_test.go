package ordering

import (
	"errors"
	"testing"
)

func TestOrderCanBeCancelledBeforeShipment(t *testing.T) {
	order, err := NewOrderFromQuote("order-001", approvedQuote(t))
	if err != nil {
		t.Fatal(err)
	}
	if err := order.Cancel(); err != nil {
		t.Fatal(err)
	}
	if order.Status() != OrderStatusCancelled {
		t.Fatalf("status = %s, want %s", order.Status(), OrderStatusCancelled)
	}
	if err := order.Cancel(); !errors.Is(err, ErrOrderNotCancellable) {
		t.Fatalf("repeated cancellation returned %v", err)
	}
}

func TestShippedOrderCannotBeCancelled(t *testing.T) {
	order, err := NewOrderFromQuote("order-001", approvedQuote(t))
	if err != nil {
		t.Fatal(err)
	}
	if err := order.MarkPaid(); err != nil {
		t.Fatal(err)
	}
	if err := order.MarkShipped(); err != nil {
		t.Fatal(err)
	}
	if err := order.Cancel(); !errors.Is(err, ErrOrderNotCancellable) {
		t.Fatalf("shipped cancellation returned %v", err)
	}
}
