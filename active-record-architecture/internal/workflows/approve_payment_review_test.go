package workflows

import (
	"testing"

	"active-record-architecture/internal/records"
)

func paymentInReview(t *testing.T) (*records.Database, *records.Order) {
	t.Helper()
	db, order := readyOrder(t)
	if _, err := CapturePayment(db, order.ID, records.PaymentOutcomeReview); err != nil {
		t.Fatalf("CapturePayment() error = %v", err)
	}
	return db, order
}

func TestApprovePaymentReviewAcceptsPayment(t *testing.T) {
	db, order := paymentInReview(t)

	resolved, err := ApprovePaymentReview(db, order.ID, "finance-1", records.PaymentReviewDecisionAccept, "manual verification passed")
	if err != nil {
		t.Fatalf("ApprovePaymentReview() error = %v", err)
	}
	if resolved.Status != records.OrderStatusReadyForFulfillment || resolved.PaymentStatus != records.PaymentStatusAccepted {
		t.Fatalf("resolved order = %#v", resolved)
	}
	payment, err := records.FindPayment(db, resolved.PaymentID)
	if err != nil {
		t.Fatalf("FindPayment() error = %v", err)
	}
	if payment.Status != records.PaymentStatusAccepted || payment.ReviewedBy != "finance-1" || payment.DecisionComment != "manual verification passed" {
		t.Fatalf("resolved payment = %#v", payment)
	}
}

func TestApprovePaymentReviewRejectsPaymentForRetry(t *testing.T) {
	db, order := paymentInReview(t)

	resolved, err := ApprovePaymentReview(db, order.ID, "finance-1", records.PaymentReviewDecisionReject, "verification failed")
	if err != nil {
		t.Fatalf("ApprovePaymentReview() error = %v", err)
	}
	if resolved.Status != records.OrderStatusReadyForPayment || resolved.PaymentStatus != records.PaymentStatusFailed {
		t.Fatalf("rejected review order = %#v", resolved)
	}
	payment, err := records.FindPayment(db, resolved.PaymentID)
	if err != nil {
		t.Fatalf("FindPayment() error = %v", err)
	}
	if payment.Status != records.PaymentStatusFailed || payment.ReviewedBy != "finance-1" {
		t.Fatalf("rejected review payment = %#v", payment)
	}
}

func TestApprovePaymentReviewValidatesBeforeMutation(t *testing.T) {
	db, order := paymentInReview(t)
	if _, err := ApprovePaymentReview(db, order.ID, "", records.PaymentReviewDecisionAccept, "missing reviewer"); err != records.ErrPaymentReviewerRequired {
		t.Fatalf("missing reviewer error = %v, want %v", err, records.ErrPaymentReviewerRequired)
	}
	if _, err := ApprovePaymentReview(db, order.ID, "finance-1", "unknown", "invalid"); err != records.ErrPaymentDecisionInvalid {
		t.Fatalf("invalid decision error = %v, want %v", err, records.ErrPaymentDecisionInvalid)
	}
	if _, err := ApprovePaymentReview(db, order.ID, "finance-1", records.PaymentReviewDecisionAccept, "accepted"); err != nil {
		t.Fatalf("valid review error = %v", err)
	}
	if _, err := ApprovePaymentReview(db, order.ID, "finance-1", records.PaymentReviewDecisionReject, "again"); err != records.ErrOrderNotInPaymentReview {
		t.Fatalf("repeated review error = %v, want %v", err, records.ErrOrderNotInPaymentReview)
	}
}

func TestApprovePaymentReviewRejectsMissingPayment(t *testing.T) {
	db, order := paymentInReview(t)
	order, err := records.FindOrder(db, order.ID)
	if err != nil {
		t.Fatalf("FindOrder() error = %v", err)
	}
	order.PaymentID = "payment-404"
	if err := order.Save(); err != nil {
		t.Fatalf("Order.Save() error = %v", err)
	}
	if _, err := ApprovePaymentReview(db, order.ID, "finance-1", records.PaymentReviewDecisionAccept, "missing payment"); err != records.ErrPaymentReviewMissing {
		t.Fatalf("missing payment error = %v, want %v", err, records.ErrPaymentReviewMissing)
	}
}
