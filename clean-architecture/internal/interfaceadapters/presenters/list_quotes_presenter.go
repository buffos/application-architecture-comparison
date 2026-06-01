package presenters

import (
	"fmt"

	"clean-architecture/internal/usecases"
)

type QuoteListItemViewModel struct {
	QuoteID    string
	CustomerID string
	Status     string
	Lines      int
}

type ListQuotesViewModel struct {
	Message string
	Status  string
	Count   int
	Quotes  []QuoteListItemViewModel
}

type ListQuotesPresenter struct {
	viewModel ListQuotesViewModel
}

func NewListQuotesPresenter() *ListQuotesPresenter {
	return &ListQuotesPresenter{}
}

func (p *ListQuotesPresenter) Present(output usecases.ListQuotesOutput) error {
	items := make([]QuoteListItemViewModel, 0, len(output.Quotes))
	for _, quote := range output.Quotes {
		items = append(items, QuoteListItemViewModel{
			QuoteID:    quote.QuoteID,
			CustomerID: quote.CustomerID,
			Status:     quote.Status,
			Lines:      quote.Lines,
		})
	}

	p.viewModel = ListQuotesViewModel{
		Message: fmt.Sprintf("listed quotes: status=%s count=%d", output.Status, output.Count),
		Status:  output.Status,
		Count:   output.Count,
		Quotes:  items,
	}

	return nil
}

func (p *ListQuotesPresenter) ViewModel() ListQuotesViewModel {
	return p.viewModel
}
