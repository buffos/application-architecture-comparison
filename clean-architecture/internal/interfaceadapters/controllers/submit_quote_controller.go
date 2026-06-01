package controllers

import "clean-architecture/internal/usecases"

type SubmitQuoteController struct {
	useCase usecases.SubmitQuoteInputBoundary
}

func NewSubmitQuoteController(useCase usecases.SubmitQuoteInputBoundary) SubmitQuoteController {
	return SubmitQuoteController{useCase: useCase}
}

func (c SubmitQuoteController) Handle(quoteID string) error {
	return c.useCase.Execute(usecases.SubmitQuoteInput{
		QuoteID: quoteID,
	})
}
