package scripts

import (
	"testing"

	"transaction-script-architecture/internal/data"
)

func TestCapturePaymentAcceptsReadyOrder(t *testing.T) {
	store := data.NewStore()
	store.Orders["order-001"] = data.Order{
		ID:            "order-001",
		Status:        data.OrderStatusReadyForPayment,
		PaymentStatus: data.PaymentStatusPending,
		Total:         30000,
	}

	got, err := CapturePayment(store, "order-001", PaymentOutcomeAccept)
	if err != nil {
		t.Fatalf("CapturePayment() error = %v", err)
	}
	if got.Status != data.OrderStatusReadyForFulfillment {
		t.Fatalf("order status = %q, want %q", got.Status, data.OrderStatusReadyForFulfillment)
	}
	if got.PaymentStatus != data.PaymentStatusAccepted {
		t.Fatalf("payment status = %q, want %q", got.PaymentStatus, data.PaymentStatusAccepted)
	}
	if got.PaymentID == "" {
		t.Fatal("payment ID is empty")
	}
	if store.Payments[got.PaymentID].Amount != 30000 {
		t.Fatalf("payment amount = %d, want 30000", store.Payments[got.PaymentID].Amount)
	}
}

func TestCapturePaymentCanLeaveOrderInManualReview(t *testing.T) {
	store := data.NewStore()
	store.Orders["order-001"] = data.Order{
		ID:     "order-001",
		Status: data.OrderStatusReadyForPayment,
		Total:  30000,
	}

	got, err := CapturePayment(store, "order-001", PaymentOutcomeReview)
	if err != nil {
		t.Fatalf("CapturePayment() error = %v", err)
	}
	if got.Status != data.OrderStatusPaymentReview {
		t.Fatalf("order status = %q, want %q", got.Status, data.OrderStatusPaymentReview)
	}
	if got.PaymentStatus != data.PaymentStatusManualReview {
		t.Fatalf("payment status = %q, want %q", got.PaymentStatus, data.PaymentStatusManualReview)
	}
}

func TestCapturePaymentFailureLeavesOrderRetryable(t *testing.T) {
	store := data.NewStore()
	store.Orders["order-001"] = data.Order{ID: "order-001", Status: data.OrderStatusReadyForPayment}

	got, err := CapturePayment(store, "order-001", PaymentOutcomeFail)
	if err != nil {
		t.Fatalf("CapturePayment() error = %v", err)
	}
	if got.Status != data.OrderStatusReadyForPayment || got.PaymentStatus != data.PaymentStatusFailed {
		t.Fatalf("order after failure = %#v, want retryable failed payment", got)
	}
}

func TestCapturePaymentRejectsInvalidStateAndOutcome(t *testing.T) {
	store := data.NewStore()
	store.Orders["order-001"] = data.Order{ID: "order-001", Status: data.OrderStatusShipped}

	if _, err := CapturePayment(store, "order-001", PaymentOutcomeAccept); err != ErrOrderNotPayable {
		t.Fatalf("state error = %v, want %v", err, ErrOrderNotPayable)
	}

	store.Orders["order-002"] = data.Order{ID: "order-002", Status: data.OrderStatusReadyForPayment}
	if _, err := CapturePayment(store, "order-002", "unknown"); err != ErrPaymentOutcomeInvalid {
		t.Fatalf("outcome error = %v, want %v", err, ErrPaymentOutcomeInvalid)
	}
}
