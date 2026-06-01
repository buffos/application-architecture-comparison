package presenters

import (
	"fmt"

	"clean-architecture/internal/usecases"
)

type CancelOrderViewModel struct {
	Message string
	OrderID string
	Status  string
	Lines   int
}

type CancelOrderPresenter struct {
	viewModel CancelOrderViewModel
}

func NewCancelOrderPresenter() *CancelOrderPresenter {
	return &CancelOrderPresenter{}
}

func (p *CancelOrderPresenter) Present(output usecases.CancelOrderOutput) error {
	p.viewModel = CancelOrderViewModel{
		Message: fmt.Sprintf("cancelled order: id=%s lines=%d status=%s", output.OrderID, output.Lines, output.Status),
		OrderID: output.OrderID,
		Status:  output.Status,
		Lines:   output.Lines,
	}

	return nil
}

func (p *CancelOrderPresenter) ViewModel() CancelOrderViewModel {
	return p.viewModel
}
