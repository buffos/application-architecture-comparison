package presenters

import (
	"fmt"

	"clean-architecture/internal/usecases"
)

type GetQuoteViewModel struct {
	Message    string
	QuoteID    string
	CustomerID string
	Status     string
	Lines      int
}

type GetQuotePresenter struct {
	viewModel GetQuoteViewModel
}

func NewGetQuotePresenter() *GetQuotePresenter {
	return &GetQuotePresenter{}
}

func (p *GetQuotePresenter) Present(output usecases.GetQuoteOutput) error {
	p.viewModel = GetQuoteViewModel{
		Message:    fmt.Sprintf("loaded quote: id=%s customer=%s lines=%d status=%s", output.QuoteID, output.CustomerID, output.Lines, output.Status),
		QuoteID:    output.QuoteID,
		CustomerID: output.CustomerID,
		Status:     output.Status,
		Lines:      output.Lines,
	}

	return nil
}

func (p *GetQuotePresenter) ViewModel() GetQuoteViewModel {
	return p.viewModel
}
