package application

import (
	"hexagonal-architecture/internal/core/domain"
	"hexagonal-architecture/internal/core/ports"
)

type RequestReturnUseCase struct {
	orders  ports.OrderRepository
	returns ports.ReturnRequestRepository
}

func NewRequestReturnUseCase(orders ports.OrderRepository, returns ports.ReturnRequestRepository) RequestReturnUseCase {
	return RequestReturnUseCase{
		orders:  orders,
		returns: returns,
	}
}

func (uc RequestReturnUseCase) Execute(orderID, reason string) (domain.ReturnRequest, error) {
	order, err := uc.orders.FindByID(orderID)
	if err != nil {
		return domain.ReturnRequest{}, err
	}

	request, err := domain.NewReturnRequest(order, reason)
	if err != nil {
		return domain.ReturnRequest{}, err
	}

	if err := uc.returns.Save(request); err != nil {
		return domain.ReturnRequest{}, err
	}

	return request, nil
}
