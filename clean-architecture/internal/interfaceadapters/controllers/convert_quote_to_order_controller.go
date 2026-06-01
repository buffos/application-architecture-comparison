package controllers

import "clean-architecture/internal/usecases"

type ConvertQuoteToOrderController struct {
	useCase usecases.ConvertQuoteToOrderInputBoundary
}

func NewConvertQuoteToOrderController(useCase usecases.ConvertQuoteToOrderInputBoundary) ConvertQuoteToOrderController {
	return ConvertQuoteToOrderController{useCase: useCase}
}

func (c ConvertQuoteToOrderController) Handle(quoteID string) error {
	return c.useCase.Execute(usecases.ConvertQuoteToOrderInput{
		QuoteID: quoteID,
	})
}
