package paymentsadapter

import "component-based-architecture/internal/components/payments"

// AcceptAllGateway is a development adapter that accepts every capture.
type AcceptAllGateway struct{}

func NewAcceptAllGateway() AcceptAllGateway {
	return AcceptAllGateway{}
}

func (g AcceptAllGateway) Capture(request payments.PaymentRequest) (payments.CaptureResult, error) {
	return payments.CaptureResult{}, nil
}

func (g AcceptAllGateway) Refund(request payments.RefundRequest) error { return nil }

var _ payments.Gateway = AcceptAllGateway{}
