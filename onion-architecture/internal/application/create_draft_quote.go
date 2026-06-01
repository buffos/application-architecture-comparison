package application

import "onion-architecture/internal/domain"

type CreateDraftQuoteCommand struct {
	CustomerID string
}

type CreateDraftQuoteResult struct {
	QuoteID    string
	CustomerID string
	Status     string
}

type QuoteRepository interface {
	Save(quote domain.Quote) error
}

type CustomerRepository interface {
	FindByID(id string) (domain.Customer, error)
}

type CreateDraftQuoteService struct {
	quotes    QuoteRepository
	customers CustomerRepository
}

func NewCreateDraftQuoteService(quotes QuoteRepository, customers CustomerRepository) CreateDraftQuoteService {
	return CreateDraftQuoteService{
		quotes:    quotes,
		customers: customers,
	}
}

func (s CreateDraftQuoteService) Execute(command CreateDraftQuoteCommand) (CreateDraftQuoteResult, error) {
	customer, err := s.customers.FindByID(command.CustomerID)
	if err != nil {
		return CreateDraftQuoteResult{}, err
	}

	if err := customer.EnsureActive(); err != nil {
		return CreateDraftQuoteResult{}, err
	}

	quote, err := domain.NewDraftQuote(command.CustomerID)
	if err != nil {
		return CreateDraftQuoteResult{}, err
	}

	if err := s.quotes.Save(quote); err != nil {
		return CreateDraftQuoteResult{}, err
	}

	return CreateDraftQuoteResult{
		QuoteID:    quote.ID,
		CustomerID: quote.CustomerID,
		Status:     quote.Status,
	}, nil
}
