package application

type SubmitQuoteCommand struct {
	QuoteID string
}

type SubmitQuoteResult struct {
	QuoteID    string
	Status     string
	LineCount  int
	TotalItems int
}

type SubmitQuoteService struct {
	quotes QuoteStore
}

func NewSubmitQuoteService(quotes QuoteStore) SubmitQuoteService {
	return SubmitQuoteService{
		quotes: quotes,
	}
}

func (s SubmitQuoteService) Execute(command SubmitQuoteCommand) (SubmitQuoteResult, error) {
	quote, err := s.quotes.FindByID(command.QuoteID)
	if err != nil {
		return SubmitQuoteResult{}, err
	}

	if err := quote.Submit(); err != nil {
		return SubmitQuoteResult{}, err
	}

	if err := s.quotes.Save(quote); err != nil {
		return SubmitQuoteResult{}, err
	}

	return SubmitQuoteResult{
		QuoteID:    quote.ID,
		Status:     quote.Status,
		LineCount:  len(quote.Lines),
		TotalItems: quote.TotalQuantity(),
	}, nil
}
