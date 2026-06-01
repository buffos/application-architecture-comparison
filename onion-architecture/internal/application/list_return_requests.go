package application

type ListReturnRequestsQuery struct {
	Status string
}

type ListReturnRequestsService struct {
	returns ReturnRequestFinder
}

func NewListReturnRequestsService(returns ReturnRequestFinder) ListReturnRequestsService {
	return ListReturnRequestsService{
		returns: returns,
	}
}

func (s ListReturnRequestsService) Execute(query ListReturnRequestsQuery) ([]ReturnRequestDetails, error) {
	requests, err := s.returns.ListByStatus(query.Status)
	if err != nil {
		return nil, err
	}

	result := make([]ReturnRequestDetails, 0, len(requests))
	for _, request := range requests {
		result = append(result, ReturnRequestDetails{
			ReturnRequestID: request.ID,
			OrderID:         request.OrderID,
			Status:          request.Status,
			Reason:          request.Reason,
			RequestedBy:     request.RequestedBy,
		})
	}

	return result, nil
}
