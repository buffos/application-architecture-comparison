package quotes

type GetQuoteQuery struct {
	QuoteID string
}

type QuoteDetails struct {
	QuoteID    string
	CustomerID string
	Status     string
	LineCount  int
}

type ListQuotesQuery struct {
	Status string
}

func (s Service) GetQuote(query GetQuoteQuery) (QuoteDetails, error) {
	quote, err := s.quotes.FindByID(query.QuoteID)
	if err != nil {
		return QuoteDetails{}, err
	}

	return QuoteDetails{
		QuoteID:    quote.ID,
		CustomerID: quote.CustomerID,
		Status:     quote.Status,
		LineCount:  len(quote.Lines),
	}, nil
}

func (s Service) ListQuotes(query ListQuotesQuery) ([]QuoteDetails, error) {
	quotes, err := s.quotes.ListByStatus(query.Status)
	if err != nil {
		return nil, err
	}

	list := make([]QuoteDetails, 0, len(quotes))
	for _, quote := range quotes {
		list = append(list, QuoteDetails{
			QuoteID:    quote.ID,
			CustomerID: quote.CustomerID,
			Status:     quote.Status,
			LineCount:  len(quote.Lines),
		})
	}

	return list, nil
}
