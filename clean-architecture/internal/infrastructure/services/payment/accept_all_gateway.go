package payment

import (
	"clean-architecture/internal/entities"
	"clean-architecture/internal/usecases"
)

type AcceptAllGateway struct{}

func NewAcceptAllGateway() AcceptAllGateway {
	return AcceptAllGateway{}
}

func (g AcceptAllGateway) Capture(order entities.Order) (string, error) {
	return usecases.PaymentCaptureApproved, nil
}
