package quotes

import "microkernel-architecture/internal/kernel"

type Service struct {
	quotes    Repository
	customers kernel.CustomerDirectory
}

func NewService(quotes Repository, customers kernel.CustomerDirectory) Service {
	return Service{
		quotes:    quotes,
		customers: customers,
	}
}

func (s Service) CreateDraftQuote(command kernel.CreateDraftQuoteCommand) (kernel.CreateDraftQuoteResult, error) {
	if err := s.customers.RequireActiveCustomer(command.CustomerID); err != nil {
		return kernel.CreateDraftQuoteResult{}, err
	}

	quote, err := NewDraftQuote(command.CustomerID)
	if err != nil {
		return kernel.CreateDraftQuoteResult{}, err
	}

	if err := s.quotes.Save(quote); err != nil {
		return kernel.CreateDraftQuoteResult{}, err
	}

	return kernel.CreateDraftQuoteResult{
		QuoteID:    quote.ID,
		CustomerID: quote.CustomerID,
		Status:     quote.Status,
	}, nil
}

func (s Service) GetQuote(query kernel.GetQuoteQuery) (kernel.QuoteDetails, error) {
	quote, err := s.quotes.FindByID(query.QuoteID)
	if err != nil {
		return kernel.QuoteDetails{}, err
	}

	return kernel.QuoteDetails{
		QuoteID:    quote.ID,
		CustomerID: quote.CustomerID,
		Status:     quote.Status,
	}, nil
}
