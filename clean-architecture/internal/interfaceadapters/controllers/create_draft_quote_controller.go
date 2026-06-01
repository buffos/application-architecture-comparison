package controllers

import "clean-architecture/internal/usecases"

type CreateDraftQuoteController struct {
	useCase usecases.CreateDraftQuoteInputBoundary
}

func NewCreateDraftQuoteController(useCase usecases.CreateDraftQuoteInputBoundary) CreateDraftQuoteController {
	return CreateDraftQuoteController{useCase: useCase}
}

func (c CreateDraftQuoteController) Handle(customerID string) error {
	return c.useCase.Execute(usecases.CreateDraftQuoteInput{
		CustomerID: customerID,
	})
}
