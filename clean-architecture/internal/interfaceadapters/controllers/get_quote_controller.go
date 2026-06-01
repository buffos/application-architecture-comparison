package controllers

import "clean-architecture/internal/usecases"

type GetQuoteController struct {
	useCase usecases.GetQuoteInputBoundary
}

func NewGetQuoteController(useCase usecases.GetQuoteInputBoundary) GetQuoteController {
	return GetQuoteController{useCase: useCase}
}

func (c GetQuoteController) Handle(quoteID string) error {
	return c.useCase.Execute(usecases.GetQuoteInput{
		QuoteID: quoteID,
	})
}
