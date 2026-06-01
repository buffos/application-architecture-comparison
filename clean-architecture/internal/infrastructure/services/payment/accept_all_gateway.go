package payment

import "clean-architecture/internal/entities"

type AcceptAllGateway struct{}

func NewAcceptAllGateway() AcceptAllGateway {
	return AcceptAllGateway{}
}

func (g AcceptAllGateway) Capture(order entities.Order) error {
	return nil
}
