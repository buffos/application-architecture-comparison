package controllers

import "clean-architecture/internal/usecases"

type RequestReturnController struct {
	useCase usecases.RequestReturnInputBoundary
}

func NewRequestReturnController(useCase usecases.RequestReturnInputBoundary) RequestReturnController {
	return RequestReturnController{useCase: useCase}
}

func (c RequestReturnController) Handle(orderID string, reason string) error {
	return c.useCase.Execute(usecases.RequestReturnInput{
		OrderID: orderID,
		Reason:  reason,
	})
}
