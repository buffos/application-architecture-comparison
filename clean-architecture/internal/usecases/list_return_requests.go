package usecases

import "clean-architecture/internal/entities"

type ListReturnRequestsInput struct {
	Status string
}

type ReturnRequestListItem struct {
	ReturnRequestID string
	OrderID         string
	Status          string
	Reason          string
}

type ListReturnRequestsOutput struct {
	Status   string
	Count    int
	Requests []ReturnRequestListItem
}

type ListReturnRequestsInputBoundary interface {
	Execute(input ListReturnRequestsInput) error
}

type ListReturnRequestsOutputBoundary interface {
	Present(output ListReturnRequestsOutput) error
}

type ReturnRequestLister interface {
	ListByStatus(status string) ([]entities.ReturnRequest, error)
}

type ListReturnRequestsInteractor struct {
	returns ReturnRequestLister
	output  ListReturnRequestsOutputBoundary
}

func NewListReturnRequestsInteractor(returns ReturnRequestLister, output ListReturnRequestsOutputBoundary) ListReturnRequestsInteractor {
	return ListReturnRequestsInteractor{
		returns: returns,
		output:  output,
	}
}

func (uc ListReturnRequestsInteractor) Execute(input ListReturnRequestsInput) error {
	requests, err := uc.returns.ListByStatus(input.Status)
	if err != nil {
		return err
	}

	items := make([]ReturnRequestListItem, 0, len(requests))
	for _, request := range requests {
		items = append(items, ReturnRequestListItem{
			ReturnRequestID: request.ID,
			OrderID:         request.OrderID,
			Status:          request.Status,
			Reason:          request.Reason,
		})
	}

	return uc.output.Present(ListReturnRequestsOutput{
		Status:   input.Status,
		Count:    len(items),
		Requests: items,
	})
}
