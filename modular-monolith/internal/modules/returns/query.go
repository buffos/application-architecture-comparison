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
	Lines           []ReturnLineDetails
}

type ListReturnRequestsQuery struct {
	Status string
}

type ReturnLineDetails struct {
	ProductSKU      string
	ProductCategory string
	Quantity        int
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
		Lines:           toReturnLineDetails(request.Lines),
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
			Lines:           toReturnLineDetails(request.Lines),
		})
	}

	return list, nil
}

func toReturnLineDetails(lines []ReturnRequestLine) []ReturnLineDetails {
	details := make([]ReturnLineDetails, 0, len(lines))
	for _, line := range lines {
		details = append(details, ReturnLineDetails{
			ProductSKU:      line.ProductSKU,
			ProductCategory: line.ProductCategory,
			Quantity:        line.Quantity,
		})
	}

	return details
}
