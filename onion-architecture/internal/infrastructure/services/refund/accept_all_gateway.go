package refund

import "onion-architecture/internal/domain"

type AcceptAllGateway struct{}

func NewAcceptAllGateway() AcceptAllGateway {
	return AcceptAllGateway{}
}

func (g AcceptAllGateway) Refund(order domain.Order) error {
	return nil
}
