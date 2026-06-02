package payment

import "modular-monolith/internal/modules/payments"

type AcceptAllGateway struct{}

func NewAcceptAllGateway() AcceptAllGateway {
	return AcceptAllGateway{}
}

func (g AcceptAllGateway) Capture(request payments.PaymentRequest) error {
	return nil
}

func (g AcceptAllGateway) Refund(request payments.RefundRequest) error {
	return nil
}
