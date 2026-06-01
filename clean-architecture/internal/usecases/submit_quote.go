package usecases

import "clean-architecture/internal/entities"

type SubmitQuoteInput struct {
	QuoteID string
}

type SubmitQuoteOutput struct {
	QuoteID string
	Status  string
	Lines   int
}

type SubmitQuoteInputBoundary interface {
	Execute(input SubmitQuoteInput) error
}

type SubmitQuoteOutputBoundary interface {
	Present(output SubmitQuoteOutput) error
}

type ApprovalPolicy interface {
	RequiresApproval(quote entities.Quote) (bool, error)
}

type SubmitQuoteInteractor struct {
	quotes QuoteEditor
	approval ApprovalPolicy
	output SubmitQuoteOutputBoundary
}

func NewSubmitQuoteInteractor(quotes QuoteEditor, approval ApprovalPolicy, output SubmitQuoteOutputBoundary) SubmitQuoteInteractor {
	return SubmitQuoteInteractor{
		quotes: quotes,
		approval: approval,
		output: output,
	}
}

func (uc SubmitQuoteInteractor) Execute(input SubmitQuoteInput) error {
	quote, err := uc.quotes.FindByID(input.QuoteID)
	if err != nil {
		return err
	}

	requiresApproval, err := uc.approval.RequiresApproval(quote)
	if err != nil {
		return err
	}

	if err := quote.Submit(requiresApproval); err != nil {
		return err
	}

	if err := uc.quotes.Save(quote); err != nil {
		return err
	}

	return uc.output.Present(SubmitQuoteOutput{
		QuoteID: quote.ID,
		Status:  quote.Status,
		Lines:   len(quote.Lines),
	})
}
