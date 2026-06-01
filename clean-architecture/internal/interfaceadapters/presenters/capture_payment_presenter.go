package presenters

import (
	"fmt"

	"clean-architecture/internal/usecases"
)

type CapturePaymentViewModel struct {
	Message string
	OrderID string
	Status  string
	Lines   int
}

type CapturePaymentPresenter struct {
	viewModel CapturePaymentViewModel
}

func NewCapturePaymentPresenter() *CapturePaymentPresenter {
	return &CapturePaymentPresenter{}
}

func (p *CapturePaymentPresenter) Present(output usecases.CapturePaymentOutput) error {
	p.viewModel = CapturePaymentViewModel{
		Message: fmt.Sprintf("captured payment: id=%s lines=%d status=%s", output.OrderID, output.Lines, output.Status),
		OrderID: output.OrderID,
		Status:  output.Status,
		Lines:   output.Lines,
	}

	return nil
}

func (p *CapturePaymentPresenter) ViewModel() CapturePaymentViewModel {
	return p.viewModel
}
