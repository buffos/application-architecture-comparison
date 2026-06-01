package usecases

import "clean-architecture/internal/entities"

type GetQuoteInput struct {
	QuoteID string
}

type GetQuoteOutput struct {
	QuoteID    string
	CustomerID string
	Status     string
	Lines      int
}

// GetQuoteInputBoundary is the interface the usecase implements
type GetQuoteInputBoundary interface {
	Execute(input GetQuoteInput) error
}

// GetQuoteOutputBoundary is the interface the usecase expect from the Presenter to implement
type GetQuoteOutputBoundary interface {
	Present(output GetQuoteOutput) error
}

type QuoteReader interface {
	FindByID(id string) (entities.Quote, error)
}

type GetQuoteInteractor struct {
	quotes QuoteReader
	output GetQuoteOutputBoundary
}

func NewGetQuoteInteractor(quotes QuoteReader, output GetQuoteOutputBoundary) GetQuoteInteractor {
	return GetQuoteInteractor{
		quotes: quotes,
		output: output,
	}
}

func (uc GetQuoteInteractor) Execute(input GetQuoteInput) error {
	quote, err := uc.quotes.FindByID(input.QuoteID)
	if err != nil {
		return err
	}

	return uc.output.Present(GetQuoteOutput{
		QuoteID:    quote.ID,
		CustomerID: quote.CustomerID,
		Status:     quote.Status,
		Lines:      len(quote.Lines),
	})
}
