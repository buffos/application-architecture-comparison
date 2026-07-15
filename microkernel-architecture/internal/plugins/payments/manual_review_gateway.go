package payments

import "microkernel-architecture/internal/kernel"

type ManualReviewGateway struct{}

func NewManualReviewGateway() ManualReviewGateway {
	return ManualReviewGateway{}
}

func (g ManualReviewGateway) Capture(orderID string, amount int) (kernel.PaymentCaptureResult, error) {
	return kernel.PaymentCaptureResult{Outcome: kernel.PaymentCaptureOutcomeReview}, nil
}

func (g ManualReviewGateway) Refund(orderID string, amount int) error {
	return nil
}
