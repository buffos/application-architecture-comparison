package presenters

import (
	"fmt"

	"clean-architecture/internal/usecases"
)

type SubmitQuoteViewModel struct {
	Message string
	QuoteID string
	Status  string
	Lines   int
}

type SubmitQuotePresenter struct {
	viewModel SubmitQuoteViewModel
}

func NewSubmitQuotePresenter() *SubmitQuotePresenter {
	return &SubmitQuotePresenter{}
}

func (p *SubmitQuotePresenter) Present(output usecases.SubmitQuoteOutput) error {
	p.viewModel = SubmitQuoteViewModel{
		Message: fmt.Sprintf("submitted quote: id=%s lines=%d status=%s", output.QuoteID, output.Lines, output.Status),
		QuoteID: output.QuoteID,
		Status:  output.Status,
		Lines:   output.Lines,
	}

	return nil
}

func (p *SubmitQuotePresenter) ViewModel() SubmitQuoteViewModel {
	return p.viewModel
}
