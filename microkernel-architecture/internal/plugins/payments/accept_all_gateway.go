package payments

import "microkernel-architecture/internal/kernel"

type AcceptAllGateway struct{}

func NewAcceptAllGateway() AcceptAllGateway {
	return AcceptAllGateway{}
}

func (g AcceptAllGateway) Capture(orderID string, amount int) (kernel.PaymentCaptureResult, error) {
	return kernel.PaymentCaptureResult{Outcome: kernel.PaymentCaptureOutcomeApproved}, nil
}

func (g AcceptAllGateway) Refund(orderID string, amount int) error {
	return nil
}
