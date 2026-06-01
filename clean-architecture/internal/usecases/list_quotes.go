package usecases

import "clean-architecture/internal/entities"

type ListQuotesInput struct {
	Status string
}

type QuoteListItem struct {
	QuoteID    string
	CustomerID string
	Status     string
	Lines      int
}

type ListQuotesOutput struct {
	Status string
	Count  int
	Quotes []QuoteListItem
}

type ListQuotesInputBoundary interface {
	Execute(input ListQuotesInput) error
}

type ListQuotesOutputBoundary interface {
	Present(output ListQuotesOutput) error
}

type QuoteLister interface {
	ListByStatus(status string) ([]entities.Quote, error)
}

type ListQuotesInteractor struct {
	quotes QuoteLister
	output ListQuotesOutputBoundary
}

func NewListQuotesInteractor(quotes QuoteLister, output ListQuotesOutputBoundary) ListQuotesInteractor {
	return ListQuotesInteractor{
		quotes: quotes,
		output: output,
	}
}

func (uc ListQuotesInteractor) Execute(input ListQuotesInput) error {
	quotes, err := uc.quotes.ListByStatus(input.Status)
	if err != nil {
		return err
	}

	items := make([]QuoteListItem, 0, len(quotes))
	for _, quote := range quotes {
		items = append(items, QuoteListItem{
			QuoteID:    quote.ID,
			CustomerID: quote.CustomerID,
			Status:     quote.Status,
			Lines:      len(quote.Lines),
		})
	}

	return uc.output.Present(ListQuotesOutput{
		Status: input.Status,
		Count:  len(items),
		Quotes: items,
	})
}
