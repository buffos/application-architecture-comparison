package presenters

import (
	"fmt"

	"clean-architecture/internal/usecases"
)

type ConvertQuoteToOrderViewModel struct {
	Message       string
	OrderID       string
	SourceQuoteID string
	Status        string
	Lines         int
}

type ConvertQuoteToOrderPresenter struct {
	viewModel ConvertQuoteToOrderViewModel
}

func NewConvertQuoteToOrderPresenter() *ConvertQuoteToOrderPresenter {
	return &ConvertQuoteToOrderPresenter{}
}

func (p *ConvertQuoteToOrderPresenter) Present(output usecases.ConvertQuoteToOrderOutput) error {
	p.viewModel = ConvertQuoteToOrderViewModel{
		Message:       fmt.Sprintf("converted order: id=%s sourceQuote=%s lines=%d status=%s", output.OrderID, output.SourceQuoteID, output.Lines, output.Status),
		OrderID:       output.OrderID,
		SourceQuoteID: output.SourceQuoteID,
		Status:        output.Status,
		Lines:         output.Lines,
	}

	return nil
}

func (p *ConvertQuoteToOrderPresenter) ViewModel() ConvertQuoteToOrderViewModel {
	return p.viewModel
}
