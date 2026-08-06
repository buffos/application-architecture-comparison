package paymentsadapter

import "component-based-architecture/internal/components/payments"

// ManualReviewGateway models a gateway that requires an operator to approve
// a capture before the order becomes paid.
type ManualReviewGateway struct{}

func NewManualReviewGateway() ManualReviewGateway { return ManualReviewGateway{} }

func (g ManualReviewGateway) Capture(request payments.PaymentRequest) (payments.CaptureResult, error) {
	return payments.CaptureResult{Outcome: payments.CaptureOutcomeReview}, nil
}

func (g ManualReviewGateway) Refund(request payments.RefundRequest) error { return nil }

var _ payments.Gateway = ManualReviewGateway{}
