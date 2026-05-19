package application

import (
	"hexagonal-architecture/internal/core/domain"
	"hexagonal-architecture/internal/core/ports"
)

type ListQuotesUseCase struct {
	repo ports.QuoteRepository
}

func NewListQuotesUseCase(repo ports.QuoteRepository) ListQuotesUseCase {
	return ListQuotesUseCase{repo: repo}
}

func (uc ListQuotesUseCase) Execute(status string) ([]domain.Quote, error) {
	return uc.repo.ListByStatus(status)
}
