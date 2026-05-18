package application

import (
	"hexagonal-architecture/internal/core/domain"
	"hexagonal-architecture/internal/core/ports"
)

type CreateDraftQuoteUseCase struct {
	repo      ports.QuoteRepository
	customers ports.CustomerLookup
}

func NewCreateDraftQuoteUseCase(repo ports.QuoteRepository, customers ports.CustomerLookup) CreateDraftQuoteUseCase {
	return CreateDraftQuoteUseCase{
		repo:      repo,
		customers: customers,
	}
}

func (uc CreateDraftQuoteUseCase) Execute(customerID string) (domain.Quote, error) {
	customer, err := uc.customers.FindByID(customerID)
	if err != nil {
		return domain.Quote{}, err
	}

	if !customer.Active {
		return domain.Quote{}, domain.ErrCustomerInactive
	}

	quote, err := domain.NewDraftQuote(customerID)
	if err != nil {
		return domain.Quote{}, err
	}

	if err := uc.repo.Save(quote); err != nil {
		return domain.Quote{}, err
	}

	return quote, nil
}
