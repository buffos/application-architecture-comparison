package presenters

import (
	"fmt"

	"clean-architecture/internal/usecases"
)

type AddQuoteLineViewModel struct {
	Message     string
	QuoteID     string
	Status      string
	Lines       int
	TotalAmount int
}

type AddQuoteLinePresenter struct {
	viewModel AddQuoteLineViewModel
}

func NewAddQuoteLinePresenter() *AddQuoteLinePresenter {
	return &AddQuoteLinePresenter{}
}

func (p *AddQuoteLinePresenter) Present(output usecases.AddQuoteLineOutput) error {
	p.viewModel = AddQuoteLineViewModel{
		Message:     fmt.Sprintf("added quote line: id=%s lines=%d total=%d status=%s", output.QuoteID, output.Lines, output.TotalAmount, output.Status),
		QuoteID:     output.QuoteID,
		Status:      output.Status,
		Lines:       output.Lines,
		TotalAmount: output.TotalAmount,
	}

	return nil
}

func (p *AddQuoteLinePresenter) ViewModel() AddQuoteLineViewModel {
	return p.viewModel
}
