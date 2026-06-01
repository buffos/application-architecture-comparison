package presenters

import (
	"fmt"

	"clean-architecture/internal/usecases"
)

type ApproveQuoteViewModel struct {
	Message string
	QuoteID string
	Status  string
	Lines   int
}

type ApproveQuotePresenter struct {
	viewModel ApproveQuoteViewModel
}

func NewApproveQuotePresenter() *ApproveQuotePresenter {
	return &ApproveQuotePresenter{}
}

func (p *ApproveQuotePresenter) Present(output usecases.ApproveQuoteOutput) error {
	p.viewModel = ApproveQuoteViewModel{
		Message: fmt.Sprintf("approved quote: id=%s lines=%d status=%s", output.QuoteID, output.Lines, output.Status),
		QuoteID: output.QuoteID,
		Status:  output.Status,
		Lines:   output.Lines,
	}

	return nil
}

func (p *ApproveQuotePresenter) ViewModel() ApproveQuoteViewModel {
	return p.viewModel
}
