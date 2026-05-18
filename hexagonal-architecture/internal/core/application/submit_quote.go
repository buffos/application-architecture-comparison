package application

import (
	"hexagonal-architecture/internal/core/domain"
	"hexagonal-architecture/internal/core/ports"
)

type SubmitQuoteUseCase struct {
	quotes   ports.QuoteRepository
	approval ports.ApprovalPolicy
}

func NewSubmitQuoteUseCase(quotes ports.QuoteRepository, approval ports.ApprovalPolicy) SubmitQuoteUseCase {
	return SubmitQuoteUseCase{
		quotes:   quotes,
		approval: approval,
	}
}

func (uc SubmitQuoteUseCase) Execute(id string) (domain.Quote, error) {
	quote, err := uc.quotes.FindByID(id)
	if err != nil {
		return domain.Quote{}, err
	}

	requiresApproval, err := uc.approval.RequiresApproval(quote)
	if err != nil {
		return domain.Quote{}, err
	}

	if err := quote.Submit(requiresApproval); err != nil {
		return domain.Quote{}, err
	}

	if err := uc.quotes.Save(quote); err != nil {
		return domain.Quote{}, err
	}

	return quote, nil
}
