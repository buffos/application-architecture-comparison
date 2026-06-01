package refund

import "clean-architecture/internal/entities"

type AcceptAllGateway struct{}

func NewAcceptAllGateway() AcceptAllGateway {
	return AcceptAllGateway{}
}

func (g AcceptAllGateway) Refund(order entities.Order) error {
	return nil
}
