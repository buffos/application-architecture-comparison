package scripts

import (
	"testing"

	"transaction-script-architecture/internal/data"
)

func TestApprovePaymentReviewAcceptsPayment(t *testing.T) {
	store := data.NewStore()
	store.Orders["order-001"] = data.Order{
		ID:            "order-001",
		Status:        data.OrderStatusPaymentReview,
		PaymentID:     "payment-001",
		PaymentStatus: data.PaymentStatusManualReview,
	}
	store.Payments["payment-001"] = data.Payment{ID: "payment-001", OrderID: "order-001", Status: data.PaymentStatusManualReview}

	got, err := ApprovePaymentReview(store, "order-001", "manager-1", PaymentReviewDecisionAccept, "Reviewed by manager")
	if err != nil {
		t.Fatalf("ApprovePaymentReview() error = %v", err)
	}
	if got.Status != data.OrderStatusReadyForFulfillment || got.PaymentStatus != data.PaymentStatusAccepted {
		t.Fatalf("order = %#v, want fulfillable accepted order", got)
	}
	if payment := store.Payments["payment-001"]; payment.Status != data.PaymentStatusAccepted || payment.ReviewedBy != "manager-1" {
		t.Fatalf("payment = %#v, want accepted with reviewer", payment)
	}
}

func TestApprovePaymentReviewRejectsPaymentForRetry(t *testing.T) {
	store := data.NewStore()
	store.Orders["order-001"] = data.Order{ID: "order-001", Status: data.OrderStatusPaymentReview, PaymentID: "payment-001"}
	store.Payments["payment-001"] = data.Payment{ID: "payment-001", Status: data.PaymentStatusManualReview}

	got, err := ApprovePaymentReview(store, "order-001", "manager-1", PaymentReviewDecisionReject, "Needs a new card")
	if err != nil {
		t.Fatalf("ApprovePaymentReview() error = %v", err)
	}
	if got.Status != data.OrderStatusReadyForPayment || got.PaymentStatus != data.PaymentStatusFailed {
		t.Fatalf("order = %#v, want retryable failed order", got)
	}
}

func TestApprovePaymentReviewRejectsInvalidInput(t *testing.T) {
	store := data.NewStore()
	store.Orders["order-001"] = data.Order{ID: "order-001", Status: data.OrderStatusReadyForPayment}

	if _, err := ApprovePaymentReview(store, "order-001", "manager-1", PaymentReviewDecisionAccept, ""); err != ErrOrderNotInPaymentReview {
		t.Fatalf("state error = %v, want %v", err, ErrOrderNotInPaymentReview)
	}
	if _, err := ApprovePaymentReview(store, "order-001", "", PaymentReviewDecisionAccept, ""); err != ErrPaymentReviewerRequired {
		t.Fatalf("reviewer error = %v, want %v", err, ErrPaymentReviewerRequired)
	}
}
