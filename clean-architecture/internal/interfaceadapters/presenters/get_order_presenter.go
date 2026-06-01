package presenters

import (
	"fmt"

	"clean-architecture/internal/usecases"
)

type GetOrderViewModel struct {
	Message       string
	OrderID       string
	CustomerID    string
	SourceQuoteID string
	Status        string
	Lines         int
}

type GetOrderPresenter struct {
	viewModel GetOrderViewModel
}

func NewGetOrderPresenter() *GetOrderPresenter {
	return &GetOrderPresenter{}
}

func (p *GetOrderPresenter) Present(output usecases.GetOrderOutput) error {
	p.viewModel = GetOrderViewModel{
		Message:       fmt.Sprintf("loaded order: id=%s customer=%s lines=%d status=%s", output.OrderID, output.CustomerID, output.Lines, output.Status),
		OrderID:       output.OrderID,
		CustomerID:    output.CustomerID,
		SourceQuoteID: output.SourceQuoteID,
		Status:        output.Status,
		Lines:         output.Lines,
	}

	return nil
}

func (p *GetOrderPresenter) ViewModel() GetOrderViewModel {
	return p.viewModel
}
