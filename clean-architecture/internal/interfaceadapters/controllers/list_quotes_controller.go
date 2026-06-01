package controllers

import "clean-architecture/internal/usecases"

type ListQuotesController struct {
	useCase usecases.ListQuotesInputBoundary
}

func NewListQuotesController(useCase usecases.ListQuotesInputBoundary) ListQuotesController {
	return ListQuotesController{useCase: useCase}
}

func (c ListQuotesController) Handle(status string) error {
	return c.useCase.Execute(usecases.ListQuotesInput{
		Status: status,
	})
}
