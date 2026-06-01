package presenters

import (
	"fmt"

	"clean-architecture/internal/usecases"
)

type RejectReturnViewModel struct {
	Message         string
	ReturnRequestID string
	OrderID         string
	Status          string
}

type RejectReturnPresenter struct {
	viewModel RejectReturnViewModel
}

func NewRejectReturnPresenter() *RejectReturnPresenter {
	return &RejectReturnPresenter{}
}

func (p *RejectReturnPresenter) Present(output usecases.RejectReturnOutput) error {
	p.viewModel = RejectReturnViewModel{
		Message:         fmt.Sprintf("rejected return: id=%s order=%s status=%s", output.ReturnRequestID, output.OrderID, output.Status),
		ReturnRequestID: output.ReturnRequestID,
		OrderID:         output.OrderID,
		Status:          output.Status,
	}

	return nil
}

func (p *RejectReturnPresenter) ViewModel() RejectReturnViewModel {
	return p.viewModel
}
