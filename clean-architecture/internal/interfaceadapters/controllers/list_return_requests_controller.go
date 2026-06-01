package controllers

import "clean-architecture/internal/usecases"

type ListReturnRequestsController struct {
	useCase usecases.ListReturnRequestsInputBoundary
}

func NewListReturnRequestsController(useCase usecases.ListReturnRequestsInputBoundary) ListReturnRequestsController {
	return ListReturnRequestsController{useCase: useCase}
}

func (c ListReturnRequestsController) Handle(status string) error {
	return c.useCase.Execute(usecases.ListReturnRequestsInput{
		Status: status,
	})
}
