package application

import "onion-architecture/internal/domain"

type RequestReturnCommand struct {
	OrderID     string
	Reason      string
	RequestedBy string
}

type RequestReturnResult struct {
	ReturnRequestID string
	OrderID         string
	Status          string
}

type ReturnRequestStore interface {
	Save(request domain.ReturnRequest) error
	FindByID(id string) (domain.ReturnRequest, error)
}

type RequestReturnService struct {
	orders  OrderRepository
	returns ReturnRequestStore
	clock   Clock
}

func NewRequestReturnService(orders OrderRepository, returns ReturnRequestStore, clock Clock) RequestReturnService {
	return RequestReturnService{
		orders:  orders,
		returns: returns,
		clock:   clock,
	}
}

func (s RequestReturnService) Execute(command RequestReturnCommand) (RequestReturnResult, error) {
	order, err := s.orders.FindByID(command.OrderID)
	if err != nil {
		return RequestReturnResult{}, err
	}

	if err := order.EnsureReturnable(); err != nil {
		return RequestReturnResult{}, err
	}

	request, err := domain.NewReturnRequest(order.ID, command.Reason, s.clock.Now(), command.RequestedBy)
	if err != nil {
		return RequestReturnResult{}, err
	}

	if err := s.returns.Save(request); err != nil {
		return RequestReturnResult{}, err
	}

	return RequestReturnResult{
		ReturnRequestID: request.ID,
		OrderID:         request.OrderID,
		Status:          request.Status,
	}, nil
}
