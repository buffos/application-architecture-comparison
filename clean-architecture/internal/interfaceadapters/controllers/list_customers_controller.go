package controllers

import "clean-architecture/internal/usecases"

type ListCustomersController struct {
	useCase usecases.ListCustomersInputBoundary
}

func NewListCustomersController(useCase usecases.ListCustomersInputBoundary) ListCustomersController {
	return ListCustomersController{useCase: useCase}
}

func (c ListCustomersController) Handle(activeOnly bool) error {
	return c.useCase.Execute(usecases.ListCustomersInput{
		ActiveOnly: activeOnly,
	})
}
