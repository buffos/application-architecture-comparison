package usecases

type ApproveQuoteInput struct {
	QuoteID string
}

type ApproveQuoteOutput struct {
	QuoteID string
	Status  string
	Lines   int
}

type ApproveQuoteInputBoundary interface {
	Execute(input ApproveQuoteInput) error
}

type ApproveQuoteOutputBoundary interface {
	Present(output ApproveQuoteOutput) error
}

type ApproveQuoteInteractor struct {
	quotes QuoteEditor
	output ApproveQuoteOutputBoundary
}

func NewApproveQuoteInteractor(quotes QuoteEditor, output ApproveQuoteOutputBoundary) ApproveQuoteInteractor {
	return ApproveQuoteInteractor{
		quotes: quotes,
		output: output,
	}
}

func (uc ApproveQuoteInteractor) Execute(input ApproveQuoteInput) error {
	quote, err := uc.quotes.FindByID(input.QuoteID)
	if err != nil {
		return err
	}

	if err := quote.Approve(); err != nil {
		return err
	}

	if err := uc.quotes.Save(quote); err != nil {
		return err
	}

	return uc.output.Present(ApproveQuoteOutput{
		QuoteID: quote.ID,
		Status:  quote.Status,
		Lines:   len(quote.Lines),
	})
}
