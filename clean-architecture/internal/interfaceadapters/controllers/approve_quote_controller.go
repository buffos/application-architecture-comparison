package controllers

import "clean-architecture/internal/usecases"

type ApproveQuoteController struct {
	useCase usecases.ApproveQuoteInputBoundary
}

func NewApproveQuoteController(useCase usecases.ApproveQuoteInputBoundary) ApproveQuoteController {
	return ApproveQuoteController{useCase: useCase}
}

func (c ApproveQuoteController) Handle(quoteID string) error {
	return c.useCase.Execute(usecases.ApproveQuoteInput{
		QuoteID: quoteID,
	})
}
