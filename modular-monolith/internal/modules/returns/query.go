package returns

type GetReturnRequestQuery struct {
	ReturnRequestID string
}

type ReturnRequestDetails struct {
	ReturnRequestID string
	OrderID         string
	CustomerID      string
	Status          string
	Reason          string
	LineCount       int
	RequestedBy     string
	ReviewedBy      string
	ProcessedBy     string
}

type ListReturnRequestsQuery struct {
	Status string
}

func (s Service) GetReturnRequest(query GetReturnRequestQuery) (ReturnRequestDetails, error) {
	request, err := s.returns.FindByID(query.ReturnRequestID)
	if err != nil {
		return ReturnRequestDetails{}, err
	}

	return ReturnRequestDetails{
		ReturnRequestID: request.ID,
		OrderID:         request.OrderID,
		CustomerID:      request.CustomerID,
		Status:          request.Status,
		Reason:          request.Reason,
		LineCount:       len(request.Lines),
		RequestedBy:     request.RequestedBy,
		ReviewedBy:      request.ReviewedBy,
		ProcessedBy:     request.ProcessedBy,
	}, nil
}

func (s Service) ListReturnRequests(query ListReturnRequestsQuery) ([]ReturnRequestDetails, error) {
	requests, err := s.returns.ListByStatus(query.Status)
	if err != nil {
		return nil, err
	}

	list := make([]ReturnRequestDetails, 0, len(requests))
	for _, request := range requests {
		list = append(list, ReturnRequestDetails{
			ReturnRequestID: request.ID,
			OrderID:         request.OrderID,
			CustomerID:      request.CustomerID,
			Status:          request.Status,
			Reason:          request.Reason,
			LineCount:       len(request.Lines),
			RequestedBy:     request.RequestedBy,
			ReviewedBy:      request.ReviewedBy,
			ProcessedBy:     request.ProcessedBy,
		})
	}

	return list, nil
}
