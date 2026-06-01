package controllers

import "clean-architecture/internal/usecases"

type ListOrdersController struct {
	useCase usecases.ListOrdersInputBoundary
}

func NewListOrdersController(useCase usecases.ListOrdersInputBoundary) ListOrdersController {
	return ListOrdersController{useCase: useCase}
}

func (c ListOrdersController) Handle(status string) error {
	return c.useCase.Execute(usecases.ListOrdersInput{
		Status: status,
	})
}
