package application

import "onion-architecture/internal/domain"

type ReturnRequestFinder interface {
	FindByID(id string) (domain.ReturnRequest, error)
	ListByStatus(status string) ([]domain.ReturnRequest, error)
}

type GetReturnRequestQuery struct {
	ReturnRequestID string
}

type ReturnRequestDetails struct {
	ReturnRequestID string
	OrderID         string
	Status          string
	Reason          string
	RequestedBy     string
}

type GetReturnRequestService struct {
	returns ReturnRequestFinder
}

func NewGetReturnRequestService(returns ReturnRequestFinder) GetReturnRequestService {
	return GetReturnRequestService{
		returns: returns,
	}
}

func (s GetReturnRequestService) Execute(query GetReturnRequestQuery) (ReturnRequestDetails, error) {
	request, err := s.returns.FindByID(query.ReturnRequestID)
	if err != nil {
		return ReturnRequestDetails{}, err
	}

	return ReturnRequestDetails{
		ReturnRequestID: request.ID,
		OrderID:         request.OrderID,
		Status:          request.Status,
		Reason:          request.Reason,
		RequestedBy:     request.RequestedBy,
	}, nil
}
