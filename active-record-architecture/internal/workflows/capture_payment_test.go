package workflows

import (
	"testing"

	"active-record-architecture/internal/records"
)

func readyOrder(t *testing.T) (*records.Database, *records.Order) {
	t.Helper()
	db, quote := approvedQuote(t)
	order, err := ConvertQuoteToOrder(db, quote.ID, "sales-1")
	if err != nil {
		t.Fatalf("ConvertQuoteToOrder() error = %v", err)
	}
	return db, order
}

func TestCapturePaymentAcceptsReadyOrder(t *testing.T) {
	db, order := readyOrder(t)

	got, err := CapturePayment(db, order.ID, records.PaymentOutcomeAccept)
	if err != nil {
		t.Fatalf("CapturePayment() error = %v", err)
	}
	if got.Status != records.OrderStatusReadyForFulfillment || got.PaymentStatus != records.PaymentStatusAccepted || got.PaymentID == "" {
		t.Fatalf("order after payment = %#v", got)
	}
	payment, err := records.FindPayment(db, got.PaymentID)
	if err != nil {
		t.Fatalf("FindPayment() error = %v", err)
	}
	if payment.Amount != got.Total || payment.Status != records.PaymentStatusAccepted {
		t.Fatalf("payment = %#v, want accepted payment for total %d", payment, got.Total)
	}
}

func TestCapturePaymentCanLeaveOrderInReviewOrRetry(t *testing.T) {
	db, order := readyOrder(t)
	got, err := CapturePayment(db, order.ID, records.PaymentOutcomeReview)
	if err != nil {
		t.Fatalf("review CapturePayment() error = %v", err)
	}
	if got.Status != records.OrderStatusPaymentReview || got.PaymentStatus != records.PaymentStatusManualReview {
		t.Fatalf("review order = %#v", got)
	}

	db, order = readyOrder(t)
	got, err = CapturePayment(db, order.ID, records.PaymentOutcomeFail)
	if err != nil {
		t.Fatalf("failed CapturePayment() error = %v", err)
	}
	if got.Status != records.OrderStatusReadyForPayment || got.PaymentStatus != records.PaymentStatusFailed {
		t.Fatalf("failed order = %#v", got)
	}
}

func TestCapturePaymentRejectsInvalidStateAndOutcome(t *testing.T) {
	db, order := readyOrder(t)
	order.Status = records.OrderStatusPendingReservation
	if err := order.Save(); err != nil {
		t.Fatalf("Order.Save() error = %v", err)
	}
	if _, err := CapturePayment(db, order.ID, records.PaymentOutcomeAccept); err != records.ErrOrderNotPayable {
		t.Fatalf("state error = %v, want %v", err, records.ErrOrderNotPayable)
	}

	db, order = readyOrder(t)
	if _, err := CapturePayment(db, order.ID, "unknown"); err != records.ErrPaymentOutcomeInvalid {
		t.Fatalf("outcome error = %v, want %v", err, records.ErrPaymentOutcomeInvalid)
	}
}
