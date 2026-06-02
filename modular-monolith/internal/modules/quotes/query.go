package quotes

type GetQuoteQuery struct {
	QuoteID string
}

type QuoteDetails struct {
	QuoteID    string
	CustomerID string
	Status     string
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
	}, nil
}
