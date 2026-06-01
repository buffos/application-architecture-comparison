package quotes

type CustomerDirectory interface {
	RequireActiveCustomer(id string) error
}

type CreateDraftQuoteCommand struct {
	CustomerID string
}

type CreateDraftQuoteResult struct {
	QuoteID    string
	CustomerID string
	Status     string
}

type Service struct {
	quotes    Repository
	customers CustomerDirectory
}

func NewService(quotes Repository, customers CustomerDirectory) Service {
	return Service{
		quotes:    quotes,
		customers: customers,
	}
}

func (s Service) CreateDraftQuote(command CreateDraftQuoteCommand) (CreateDraftQuoteResult, error) {
	if err := s.customers.RequireActiveCustomer(command.CustomerID); err != nil {
		return CreateDraftQuoteResult{}, err
	}

	quote, err := NewDraftQuote(command.CustomerID)
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
