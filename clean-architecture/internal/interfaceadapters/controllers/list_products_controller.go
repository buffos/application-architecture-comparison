package controllers

import "clean-architecture/internal/usecases"

type ListProductsController struct {
	useCase usecases.ListProductsInputBoundary
}

func NewListProductsController(useCase usecases.ListProductsInputBoundary) ListProductsController {
	return ListProductsController{useCase: useCase}
}

func (c ListProductsController) Handle(category string, availableOnly bool) error {
	return c.useCase.Execute(usecases.ListProductsInput{
		Category:      category,
		AvailableOnly: availableOnly,
	})
}
