package returns

import (
	"errors"
	"sort"

	domainreturns "domain-driven-design-architecture/internal/domain/returns"
)

var ErrReturnRequestNotFound = errors.New("return request not found")

type Reader interface {
	GetReturnRequest(id string) (ReturnRequestDetails, error)
	ListReturnRequests(status domainreturns.ReturnStatus) []ReturnRequestSummary
}

type ReturnRequestDetails struct {
	ReturnRequestID string
	OrderID         string
	CustomerID      string
	Reason          string
	Status          string
	RequestedBy     string
	ReviewedBy      string
	ProcessedBy     string
	Lines           []ReturnLineDetails
}

type ReturnLineDetails struct {
	ProductSKU string
	Category   string
	Quantity   int
}

type ReturnRequestSummary struct {
	ReturnRequestID string
	OrderID         string
	Status          string
	LineCount       int
}

type InMemoryReader struct {
	requests map[string]ReturnRequestDetails
}

func NewInMemoryReader() *InMemoryReader {
	return &InMemoryReader{requests: make(map[string]ReturnRequestDetails)}
}

func (r *InMemoryReader) Save(request domainreturns.ReturnRequest) {
	lines := request.Lines()
	details := ReturnRequestDetails{
		ReturnRequestID: string(request.ID()),
		OrderID:         string(request.OrderID()),
		CustomerID:      string(request.CustomerID()),
		Reason:          request.Reason(),
		Status:          string(request.Status()),
		RequestedBy:     request.RequestedBy(),
		ReviewedBy:      request.ReviewedBy(),
		ProcessedBy:     request.ProcessedBy(),
		Lines:           make([]ReturnLineDetails, 0, len(lines)),
	}
	for _, line := range lines {
		details.Lines = append(details.Lines, ReturnLineDetails{ProductSKU: string(line.ProductSKU()), Category: string(line.ProductCategory()), Quantity: line.Quantity()})
	}
	r.requests[details.ReturnRequestID] = details
}

func (r *InMemoryReader) GetReturnRequest(id string) (ReturnRequestDetails, error) {
	details, ok := r.requests[id]
	if !ok {
		return ReturnRequestDetails{}, ErrReturnRequestNotFound
	}
	details.Lines = append([]ReturnLineDetails(nil), details.Lines...)
	return details, nil
}

func (r *InMemoryReader) ListReturnRequests(status domainreturns.ReturnStatus) []ReturnRequestSummary {
	result := make([]ReturnRequestSummary, 0, len(r.requests))
	for _, details := range r.requests {
		if status != "" && details.Status != string(status) {
			continue
		}
		result = append(result, ReturnRequestSummary{ReturnRequestID: details.ReturnRequestID, OrderID: details.OrderID, Status: details.Status, LineCount: len(details.Lines)})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ReturnRequestID < result[j].ReturnRequestID })
	return result
}

var _ Reader = (*InMemoryReader)(nil)
