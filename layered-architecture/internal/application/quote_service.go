package application

import "layered-architecture/internal/domain"

type QuoteRepository interface {
	Save(quote domain.Quote) error
	FindByID(id string) (domain.Quote, error)
}

type QuoteService struct {
	repo QuoteRepository
}

func NewQuoteService(repo QuoteRepository) QuoteService {
	return QuoteService{repo: repo}
}

func (s QuoteService) CreateDraftQuote(customerID string) (domain.Quote, error) {
	quote, err := domain.NewDraftQuote(customerID)
	if err != nil {
		return domain.Quote{}, err
	}

	if err := s.repo.Save(quote); err != nil {
		return domain.Quote{}, err
	}

	return quote, nil
}

func (s QuoteService) GetQuote(id string) (domain.Quote, error) {
	return s.repo.FindByID(id)
}
