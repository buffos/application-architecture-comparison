package presenters

import (
	"fmt"

	"clean-architecture/internal/usecases"
)

type AcceptReturnViewModel struct {
	Message         string
	ReturnRequestID string
	OrderID         string
	Status          string
}

type AcceptReturnPresenter struct {
	viewModel AcceptReturnViewModel
}

func NewAcceptReturnPresenter() *AcceptReturnPresenter {
	return &AcceptReturnPresenter{}
}

func (p *AcceptReturnPresenter) Present(output usecases.AcceptReturnOutput) error {
	p.viewModel = AcceptReturnViewModel{
		Message:         fmt.Sprintf("accepted return: id=%s order=%s status=%s", output.ReturnRequestID, output.OrderID, output.Status),
		ReturnRequestID: output.ReturnRequestID,
		OrderID:         output.OrderID,
		Status:          output.Status,
	}

	return nil
}

func (p *AcceptReturnPresenter) ViewModel() AcceptReturnViewModel {
	return p.viewModel
}
