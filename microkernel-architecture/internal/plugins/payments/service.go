package payments

import "microkernel-architecture/internal/kernel"

type Gateway interface {
	Capture(orderID string, amount int) (kernel.PaymentCaptureResult, error)
	Refund(orderID string, amount int) error
}

type Service struct {
	gateway Gateway
}

func NewService(gateway Gateway) Service {
	return Service{gateway: gateway}
}

func (s Service) Capture(orderID string, amount int) (kernel.PaymentCaptureResult, error) {
	return s.gateway.Capture(orderID, amount)
}

func (s Service) Refund(orderID string, amount int) error {
	return s.gateway.Refund(orderID, amount)
}
