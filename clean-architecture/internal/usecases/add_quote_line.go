package usecases

import "clean-architecture/internal/entities"

type AddQuoteLineInput struct {
	QuoteID  string
	SKU      string
	Quantity int
}

type AddQuoteLineOutput struct {
	QuoteID string
	Status  string
	Lines   int
}

type AddQuoteLineInputBoundary interface {
	Execute(input AddQuoteLineInput) error
}

type AddQuoteLineOutputBoundary interface {
	Present(output AddQuoteLineOutput) error
}

type QuoteEditor interface {
	FindByID(id string) (entities.Quote, error)
	Save(quote entities.Quote) error
}

type ProductGateway interface {
	FindBySKU(sku string) (entities.Product, error)
}

type AddQuoteLineInteractor struct {
	quotes   QuoteEditor
	products ProductGateway
	output   AddQuoteLineOutputBoundary
}

func NewAddQuoteLineInteractor(quotes QuoteEditor, products ProductGateway, output AddQuoteLineOutputBoundary) AddQuoteLineInteractor {
	return AddQuoteLineInteractor{
		quotes:   quotes,
		products: products,
		output:   output,
	}
}

func (uc AddQuoteLineInteractor) Execute(input AddQuoteLineInput) error {
	quote, err := uc.quotes.FindByID(input.QuoteID)
	if err != nil {
		return err
	}

	product, err := uc.products.FindBySKU(input.SKU)
	if err != nil {
		return err
	}

	if err := product.EnsureAvailable(); err != nil {
		return err
	}

	if err := quote.AddLine(product, input.Quantity); err != nil {
		return err
	}

	if err := uc.quotes.Save(quote); err != nil {
		return err
	}

	return uc.output.Present(AddQuoteLineOutput{
		QuoteID: quote.ID,
		Status:  quote.Status,
		Lines:   len(quote.Lines),
	})
}
