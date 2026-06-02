package payment

import "modular-monolith/internal/modules/payments"

type ManualReviewGateway struct{}

func NewManualReviewGateway() ManualReviewGateway {
	return ManualReviewGateway{}
}

func (g ManualReviewGateway) Capture(request payments.PaymentRequest) (payments.CaptureResult, error) {
	return payments.CaptureResult{Outcome: payments.CaptureOutcomeReview}, nil
}

func (g ManualReviewGateway) Refund(request payments.RefundRequest) error {
	return nil
}
