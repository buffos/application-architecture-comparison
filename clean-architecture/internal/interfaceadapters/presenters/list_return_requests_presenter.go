package presenters

import (
	"fmt"

	"clean-architecture/internal/usecases"
)

type ReturnRequestListItemViewModel struct {
	ReturnRequestID string
	OrderID         string
	Status          string
	Reason          string
}

type ListReturnRequestsViewModel struct {
	Message  string
	Status   string
	Count    int
	Requests []ReturnRequestListItemViewModel
}

type ListReturnRequestsPresenter struct {
	viewModel ListReturnRequestsViewModel
}

func NewListReturnRequestsPresenter() *ListReturnRequestsPresenter {
	return &ListReturnRequestsPresenter{}
}

func (p *ListReturnRequestsPresenter) Present(output usecases.ListReturnRequestsOutput) error {
	items := make([]ReturnRequestListItemViewModel, 0, len(output.Requests))
	for _, request := range output.Requests {
		items = append(items, ReturnRequestListItemViewModel{
			ReturnRequestID: request.ReturnRequestID,
			OrderID:         request.OrderID,
			Status:          request.Status,
			Reason:          request.Reason,
		})
	}

	p.viewModel = ListReturnRequestsViewModel{
		Message:  fmt.Sprintf("listed return requests: status=%s count=%d", output.Status, output.Count),
		Status:   output.Status,
		Count:    output.Count,
		Requests: items,
	}

	return nil
}

func (p *ListReturnRequestsPresenter) ViewModel() ListReturnRequestsViewModel {
	return p.viewModel
}
