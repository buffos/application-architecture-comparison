package presenters

import (
	"fmt"

	"clean-architecture/internal/usecases"
)

type CreateDraftQuoteViewModel struct {
	Message    string
	QuoteID    string
	CustomerID string
	Status     string
}

type CreateDraftQuotePresenter struct {
	viewModel CreateDraftQuoteViewModel
}

func NewCreateDraftQuotePresenter() *CreateDraftQuotePresenter {
	return &CreateDraftQuotePresenter{}
}

func (p *CreateDraftQuotePresenter) Present(output usecases.CreateDraftQuoteOutput) error {
	p.viewModel = CreateDraftQuoteViewModel{
		Message:    fmt.Sprintf("created draft quote: id=%s customer=%s status=%s", output.QuoteID, output.CustomerID, output.Status),
		QuoteID:    output.QuoteID,
		CustomerID: output.CustomerID,
		Status:     output.Status,
	}

	return nil
}

func (p *CreateDraftQuotePresenter) ViewModel() CreateDraftQuoteViewModel {
	return p.viewModel
}
