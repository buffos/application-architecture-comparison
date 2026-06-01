package usecases

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

type SubmitQuoteInteractor struct {
	quotes QuoteEditor
	output SubmitQuoteOutputBoundary
}

func NewSubmitQuoteInteractor(quotes QuoteEditor, output SubmitQuoteOutputBoundary) SubmitQuoteInteractor {
	return SubmitQuoteInteractor{
		quotes: quotes,
		output: output,
	}
}

func (uc SubmitQuoteInteractor) Execute(input SubmitQuoteInput) error {
	quote, err := uc.quotes.FindByID(input.QuoteID)
	if err != nil {
		return err
	}

	if err := quote.Submit(); err != nil {
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
