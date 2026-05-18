package application

import (
	"hexagonal-architecture/internal/core/domain"
	"hexagonal-architecture/internal/core/ports"
)

type CreateDraftQuoteUseCase struct {
	repo ports.QuoteRepository
}

func NewCreateDraftQuoteUseCase(repo ports.QuoteRepository) CreateDraftQuoteUseCase {
	return CreateDraftQuoteUseCase{repo: repo}
}

func (uc CreateDraftQuoteUseCase) Execute(customerID string) (domain.Quote, error) {
	quote, err := domain.NewDraftQuote(customerID)
	if err != nil {
		return domain.Quote{}, err
	}

	if err := uc.repo.Save(quote); err != nil {
		return domain.Quote{}, err
	}

	return quote, nil
}
