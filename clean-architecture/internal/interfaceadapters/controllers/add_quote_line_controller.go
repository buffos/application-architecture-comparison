package controllers

import "clean-architecture/internal/usecases"

type AddQuoteLineController struct {
	useCase usecases.AddQuoteLineInputBoundary
}

func NewAddQuoteLineController(useCase usecases.AddQuoteLineInputBoundary) AddQuoteLineController {
	return AddQuoteLineController{useCase: useCase}
}

func (c AddQuoteLineController) Handle(quoteID string, sku string, quantity int) error {
	return c.useCase.Execute(usecases.AddQuoteLineInput{
		QuoteID:  quoteID,
		SKU:      sku,
		Quantity: quantity,
	})
}
