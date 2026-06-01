package application

type ApproveQuoteCommand struct {
	QuoteID string
}

type ApproveQuoteResult struct {
	QuoteID    string
	Status     string
	LineCount  int
	TotalItems int
}

type ApproveQuoteService struct {
	quotes QuoteStore
}

func NewApproveQuoteService(quotes QuoteStore) ApproveQuoteService {
	return ApproveQuoteService{
		quotes: quotes,
	}
}

func (s ApproveQuoteService) Execute(command ApproveQuoteCommand) (ApproveQuoteResult, error) {
	quote, err := s.quotes.FindByID(command.QuoteID)
	if err != nil {
		return ApproveQuoteResult{}, err
	}

	if err := quote.Approve(); err != nil {
		return ApproveQuoteResult{}, err
	}

	if err := s.quotes.Save(quote); err != nil {
		return ApproveQuoteResult{}, err
	}

	return ApproveQuoteResult{
		QuoteID:    quote.ID,
		Status:     quote.Status,
		LineCount:  len(quote.Lines),
		TotalItems: quote.TotalQuantity(),
	}, nil
}
