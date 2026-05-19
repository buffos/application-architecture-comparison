package application

import (
	"hexagonal-architecture/internal/core/domain"
	"hexagonal-architecture/internal/core/ports"
)

type RequestReturnUseCase struct {
	orders  ports.OrderRepository
	returns ports.ReturnRequestRepository
	clock   ports.Clock
}

func NewRequestReturnUseCase(orders ports.OrderRepository, returns ports.ReturnRequestRepository, clock ports.Clock) RequestReturnUseCase {
	return RequestReturnUseCase{
		orders:  orders,
		returns: returns,
		clock:   clock,
	}
}

func (uc RequestReturnUseCase) Execute(orderID, reason, requestedBy string) (domain.ReturnRequest, error) {
	order, err := uc.orders.FindByID(orderID)
	if err != nil {
		return domain.ReturnRequest{}, err
	}

	request, err := domain.NewReturnRequest(order, reason, requestedBy, uc.clock.Now())
	if err != nil {
		return domain.ReturnRequest{}, err
	}

	if err := uc.returns.Save(request); err != nil {
		return domain.ReturnRequest{}, err
	}

	return request, nil
}
