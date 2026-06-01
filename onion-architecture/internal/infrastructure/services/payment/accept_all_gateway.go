package payment

import "onion-architecture/internal/domain"

type AcceptAllGateway struct{}

func NewAcceptAllGateway() AcceptAllGateway {
	return AcceptAllGateway{}
}

func (g AcceptAllGateway) Capture(order domain.Order) (string, error) {
	return "Approved", nil
}
