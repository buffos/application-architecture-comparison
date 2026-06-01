package usecases

import "clean-architecture/internal/entities"

type CreateDraftQuoteInput struct {
	CustomerID string
}

type CreateDraftQuoteOutput struct {
	QuoteID    string
	CustomerID string
	Status     string
}

type CreateDraftQuoteInputBoundary interface {
	Execute(input CreateDraftQuoteInput) error
}

type CreateDraftQuoteOutputBoundary interface {
	Present(output CreateDraftQuoteOutput) error
}

type QuoteGateway interface {
	Save(quote entities.Quote) error
}

type CustomerGateway interface {
	FindByID(id string) (entities.Customer, error)
}

type CreateDraftQuoteInteractor struct {
	quotes    QuoteGateway
	customers CustomerGateway
	output    CreateDraftQuoteOutputBoundary
}

func NewCreateDraftQuoteInteractor(quotes QuoteGateway, customers CustomerGateway, output CreateDraftQuoteOutputBoundary) CreateDraftQuoteInteractor {
	return CreateDraftQuoteInteractor{
		quotes:    quotes,
		customers: customers,
		output:    output,
	}
}

func (uc CreateDraftQuoteInteractor) Execute(input CreateDraftQuoteInput) error {
	customer, err := uc.customers.FindByID(input.CustomerID)
	if err != nil {
		return err
	}

	if err := customer.EnsureActive(); err != nil {
		return err
	}

	quote, err := entities.NewDraftQuote(input.CustomerID)
	if err != nil {
		return err
	}

	if err := uc.quotes.Save(quote); err != nil {
		return err
	}

	return uc.output.Present(CreateDraftQuoteOutput{
		QuoteID:    quote.ID,
		CustomerID: quote.CustomerID,
		Status:     quote.Status,
	})
}
