package application

type ListQuotesQuery struct {
	Status string
}

type ListQuotesService struct {
	quotes QuoteFinder
}

func NewListQuotesService(quotes QuoteFinder) ListQuotesService {
	return ListQuotesService{
		quotes: quotes,
	}
}

func (s ListQuotesService) Execute(query ListQuotesQuery) ([]QuoteDetails, error) {
	quotes, err := s.quotes.ListByStatus(query.Status)
	if err != nil {
		return nil, err
	}

	result := make([]QuoteDetails, 0, len(quotes))
	for _, quote := range quotes {
		result = append(result, QuoteDetails{
			QuoteID:    quote.ID,
			CustomerID: quote.CustomerID,
			Status:     quote.Status,
			LineCount:  len(quote.Lines),
		})
	}

	return result, nil
}
