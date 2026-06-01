package application

import "onion-architecture/internal/domain"

type SubmitQuoteCommand struct {
	QuoteID string
}

type SubmitQuoteResult struct {
	QuoteID    string
	Status     string
	LineCount  int
	TotalItems int
}

type SubmitQuoteService struct {
	quotes QuoteStore
	policy ApprovalPolicy
}

type ApprovalPolicy interface {
	RequiresApproval(quote domain.Quote) (bool, error)
}

func NewSubmitQuoteService(quotes QuoteStore, policy ApprovalPolicy) SubmitQuoteService {
	return SubmitQuoteService{
		quotes: quotes,
		policy: policy,
	}
}

func (s SubmitQuoteService) Execute(command SubmitQuoteCommand) (SubmitQuoteResult, error) {
	quote, err := s.quotes.FindByID(command.QuoteID)
	if err != nil {
		return SubmitQuoteResult{}, err
	}

	requiresApproval, err := s.policy.RequiresApproval(quote)
	if err != nil {
		return SubmitQuoteResult{}, err
	}

	if err := quote.Submit(requiresApproval); err != nil {
		return SubmitQuoteResult{}, err
	}

	if err := s.quotes.Save(quote); err != nil {
		return SubmitQuoteResult{}, err
	}

	return SubmitQuoteResult{
		QuoteID:    quote.ID,
		Status:     quote.Status,
		LineCount:  len(quote.Lines),
		TotalItems: quote.TotalQuantity(),
	}, nil
}
