package payment

import "hexagonal-architecture/internal/core/domain"

type AcceptAllGateway struct{}

func NewAcceptAllGateway() AcceptAllGateway {
	return AcceptAllGateway{}
}

func (AcceptAllGateway) Capture(order domain.Order) (bool, error) {
	return true, nil
}
