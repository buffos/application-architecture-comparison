package application

import (
	"hexagonal-architecture/internal/core/domain"
	"hexagonal-architecture/internal/core/ports"
)

type ApproveQuoteUseCase struct {
	quotes ports.QuoteRepository
}

func NewApproveQuoteUseCase(quotes ports.QuoteRepository) ApproveQuoteUseCase {
	return ApproveQuoteUseCase{quotes: quotes}
}

func (uc ApproveQuoteUseCase) Execute(id string) (domain.Quote, error) {
	quote, err := uc.quotes.FindByID(id)
	if err != nil {
		return domain.Quote{}, err
	}

	if err := quote.Approve(); err != nil {
		return domain.Quote{}, err
	}

	if err := uc.quotes.Save(quote); err != nil {
		return domain.Quote{}, err
	}

	return quote, nil
}
