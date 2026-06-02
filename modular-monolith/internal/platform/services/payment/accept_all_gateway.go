package payment

import "modular-monolith/internal/modules/payments"

type AcceptAllGateway struct{}

func NewAcceptAllGateway() AcceptAllGateway {
	return AcceptAllGateway{}
}

func (g AcceptAllGateway) Capture(request payments.PaymentRequest) (payments.CaptureResult, error) {
	return payments.CaptureResult{Outcome: payments.CaptureOutcomeApproved}, nil
}

func (g AcceptAllGateway) Refund(request payments.RefundRequest) error {
	return nil
}
