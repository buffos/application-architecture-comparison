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
}

type GetQuotePresenter struct {
	viewModel GetQuoteViewModel
}

func NewGetQuotePresenter() *GetQuotePresenter {
	return &GetQuotePresenter{}
}

func (p *GetQuotePresenter) Present(output usecases.GetQuoteOutput) error {
	p.viewModel = GetQuoteViewModel{
		Message:    fmt.Sprintf("loaded quote: id=%s customer=%s status=%s", output.QuoteID, output.CustomerID, output.Status),
		QuoteID:    output.QuoteID,
		CustomerID: output.CustomerID,
		Status:     output.Status,
	}

	return nil
}

func (p *GetQuotePresenter) ViewModel() GetQuoteViewModel {
	return p.viewModel
}
