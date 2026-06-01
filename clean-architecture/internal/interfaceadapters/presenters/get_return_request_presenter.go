package presenters

import (
	"fmt"

	"clean-architecture/internal/usecases"
)

type GetReturnRequestViewModel struct {
	Message         string
	ReturnRequestID string
	OrderID         string
	Status          string
	Reason          string
	RequestedBy     string
	ReviewedBy      string
	ProcessedBy     string
}

type GetReturnRequestPresenter struct {
	viewModel GetReturnRequestViewModel
}

func NewGetReturnRequestPresenter() *GetReturnRequestPresenter {
	return &GetReturnRequestPresenter{}
}

func (p *GetReturnRequestPresenter) Present(output usecases.GetReturnRequestOutput) error {
	p.viewModel = GetReturnRequestViewModel{
		Message:         fmt.Sprintf("loaded return request: id=%s order=%s status=%s", output.ReturnRequestID, output.OrderID, output.Status),
		ReturnRequestID: output.ReturnRequestID,
		OrderID:         output.OrderID,
		Status:          output.Status,
		Reason:          output.Reason,
		RequestedBy:     output.RequestedBy,
		ReviewedBy:      output.ReviewedBy,
		ProcessedBy:     output.ProcessedBy,
	}

	return nil
}

func (p *GetReturnRequestPresenter) ViewModel() GetReturnRequestViewModel {
	return p.viewModel
}
