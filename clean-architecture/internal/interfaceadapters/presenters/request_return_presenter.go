package presenters

import (
	"fmt"

	"clean-architecture/internal/usecases"
)

type RequestReturnViewModel struct {
	Message         string
	ReturnRequestID string
	OrderID         string
	Status          string
}

type RequestReturnPresenter struct {
	viewModel RequestReturnViewModel
}

func NewRequestReturnPresenter() *RequestReturnPresenter {
	return &RequestReturnPresenter{}
}

func (p *RequestReturnPresenter) Present(output usecases.RequestReturnOutput) error {
	p.viewModel = RequestReturnViewModel{
		Message:         fmt.Sprintf("requested return: id=%s order=%s status=%s", output.ReturnRequestID, output.OrderID, output.Status),
		ReturnRequestID: output.ReturnRequestID,
		OrderID:         output.OrderID,
		Status:          output.Status,
	}

	return nil
}

func (p *RequestReturnPresenter) ViewModel() RequestReturnViewModel {
	return p.viewModel
}
