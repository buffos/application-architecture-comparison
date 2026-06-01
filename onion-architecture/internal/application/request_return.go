package application

import "onion-architecture/internal/domain"

type RequestReturnCommand struct {
	OrderID string
	Reason  string
}

type RequestReturnResult struct {
	ReturnRequestID string
	OrderID         string
	Status          string
}

type ReturnRequestStore interface {
	Save(request domain.ReturnRequest) error
}

type RefundGateway interface {
	Refund(order domain.Order) error
}

type RequestReturnService struct {
	orders  OrderRepository
	returns ReturnRequestStore
	refunds RefundGateway
}

func NewRequestReturnService(orders OrderRepository, returns ReturnRequestStore, refunds RefundGateway) RequestReturnService {
	return RequestReturnService{
		orders:  orders,
		returns: returns,
		refunds: refunds,
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

	if err := s.refunds.Refund(order); err != nil {
		return RequestReturnResult{}, err
	}

	request, err := domain.NewReturnRequest(order.ID, command.Reason)
	if err != nil {
		return RequestReturnResult{}, err
	}

	request.Status = domain.ReturnRequestStatusRefunded

	if err := s.returns.Save(request); err != nil {
		return RequestReturnResult{}, err
	}

	return RequestReturnResult{
		ReturnRequestID: request.ID,
		OrderID:         request.OrderID,
		Status:          request.Status,
	}, nil
}
