package payment

import (
	"hexagonal-architecture/internal/core/domain"
	"hexagonal-architecture/internal/core/ports"
)

type AcceptAllGateway struct{}

func NewAcceptAllGateway() AcceptAllGateway {
	return AcceptAllGateway{}
}

func (AcceptAllGateway) Capture(order domain.Order) (ports.PaymentResult, error) {
	return ports.PaymentResultAccepted, nil
}
