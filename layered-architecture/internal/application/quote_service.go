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

func (s QuoteService) AddQuoteLine(id string, productName string, quantity int) (domain.Quote, error) {
	quote, err := s.repo.FindByID(id)
	if err != nil {
		return domain.Quote{}, err
	}

	if err := quote.AddLine(productName, quantity); err != nil {
		return domain.Quote{}, err
	}

	if err := s.repo.Save(quote); err != nil {
		return domain.Quote{}, err
	}

	return quote, nil
}

func (s QuoteService) SubmitQuote(id string) (domain.Quote, error) {
	quote, err := s.repo.FindByID(id)
	if err != nil {
		return domain.Quote{}, err
	}

	if err := quote.Submit(); err != nil {
		return domain.Quote{}, err
	}

	if err := s.repo.Save(quote); err != nil {
		return domain.Quote{}, err
	}

	return quote, nil
}
