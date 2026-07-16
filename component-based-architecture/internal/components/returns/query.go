package returns

import "errors"

var ErrReturnRequestNotFound = errors.New("return request not found")

// Reader is the public read contract provided by Returns. It exposes business
// views rather than the component's private return-request map.
type Reader interface {
	GetReturnRequest(query GetReturnRequestQuery) (ReturnRequestDetails, error)
	ListReturnRequests(query ListReturnRequestsQuery) []ReturnRequestSummary
}

type GetReturnRequestQuery struct {
	ReturnRequestID string
}

type ListReturnRequestsQuery struct {
	Status string
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
	ReviewNote      string
}

type ReturnRequestSummary struct {
	ReturnRequestID string
	OrderID         string
	CustomerID      string
	Status          string
	LineCount       int
}

func (c *Component) GetReturnRequest(query GetReturnRequestQuery) (ReturnRequestDetails, error) {
	request, ok := c.requests[query.ReturnRequestID]
	if !ok {
		return ReturnRequestDetails{}, ErrReturnRequestNotFound
	}
	return returnRequestDetails(request), nil
}

func (c *Component) ListReturnRequests(query ListReturnRequestsQuery) []ReturnRequestSummary {
	requests := make([]ReturnRequestSummary, 0, len(c.requests))
	for _, request := range c.requests {
		if query.Status != "" && request.Status != query.Status {
			continue
		}
		requests = append(requests, ReturnRequestSummary{
			ReturnRequestID: request.ID,
			OrderID:         request.OrderID,
			CustomerID:      request.CustomerID,
			Status:          request.Status,
			LineCount:       request.LineCount,
		})
	}
	return requests
}

func returnRequestDetails(request ReturnRequest) ReturnRequestDetails {
	return ReturnRequestDetails{
		ReturnRequestID: request.ID,
		OrderID:         request.OrderID,
		CustomerID:      request.CustomerID,
		Status:          request.Status,
		Reason:          request.Reason,
		LineCount:       request.LineCount,
		RequestedBy:     request.RequestedBy,
		ReviewedBy:      request.ReviewedBy,
		ProcessedBy:     request.ProcessedBy,
		ReviewNote:      request.ReviewNote,
	}
}
