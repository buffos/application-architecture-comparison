package refund

import "hexagonal-architecture/internal/core/domain"

type AcceptAllGateway struct{}

func NewAcceptAllGateway() AcceptAllGateway {
	return AcceptAllGateway{}
}

func (AcceptAllGateway) Refund(request domain.ReturnRequest) (bool, error) {
	return true, nil
}
