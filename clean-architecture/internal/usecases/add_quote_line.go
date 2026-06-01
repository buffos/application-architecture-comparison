package usecases

import "clean-architecture/internal/entities"

type AddQuoteLineInput struct {
	QuoteID  string
	SKU      string
	Quantity int
}

type AddQuoteLineOutput struct {
	QuoteID     string
	Status      string
	Lines       int
	TotalAmount int
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

type PricingPolicy interface {
	AdjustUnitPrice(product entities.Product, quantity int) (int, error)
}

type AddQuoteLineInteractor struct {
	quotes   QuoteEditor
	products ProductGateway
	pricing  PricingPolicy
	output   AddQuoteLineOutputBoundary
}

func NewAddQuoteLineInteractor(quotes QuoteEditor, products ProductGateway, pricing PricingPolicy, output AddQuoteLineOutputBoundary) AddQuoteLineInteractor {
	return AddQuoteLineInteractor{
		quotes:   quotes,
		products: products,
		pricing:  pricing,
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

	adjustedUnitPrice, err := uc.pricing.AdjustUnitPrice(product, input.Quantity)
	if err != nil {
		return err
	}

	pricedProduct := product
	pricedProduct.BasePrice = adjustedUnitPrice

	if err := quote.AddLine(pricedProduct, input.Quantity); err != nil {
		return err
	}

	if err := uc.quotes.Save(quote); err != nil {
		return err
	}

	return uc.output.Present(AddQuoteLineOutput{
		QuoteID:     quote.ID,
		Status:      quote.Status,
		Lines:       len(quote.Lines),
		TotalAmount: quoteTotalAmount(quote),
	})
}

func quoteTotalAmount(quote entities.Quote) int {
	total := 0
	for _, line := range quote.Lines {
		total += line.LineTotal
	}

	return total
}
