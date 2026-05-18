package application

import (
	"hexagonal-architecture/internal/core/domain"
	"hexagonal-architecture/internal/core/ports"
)

type GetQuoteUseCase struct {
	repo ports.QuoteRepository
}

func NewGetQuoteUseCase(repo ports.QuoteRepository) GetQuoteUseCase {
	return GetQuoteUseCase{repo: repo}
}

func (uc GetQuoteUseCase) Execute(id string) (domain.Quote, error) {
	return uc.repo.FindByID(id)
}
