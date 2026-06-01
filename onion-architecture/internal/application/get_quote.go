package application

import "onion-architecture/internal/domain"

type GetQuoteQuery struct {
	QuoteID string
}

type QuoteDetails struct {
	QuoteID    string
	CustomerID string
	Status     string
}

type QuoteFinder interface {
	FindByID(id string) (domain.Quote, error)
}

type GetQuoteService struct {
	quotes QuoteFinder
}

func NewGetQuoteService(quotes QuoteFinder) GetQuoteService {
	return GetQuoteService{
		quotes: quotes,
	}
}

func (s GetQuoteService) Execute(query GetQuoteQuery) (QuoteDetails, error) {
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
