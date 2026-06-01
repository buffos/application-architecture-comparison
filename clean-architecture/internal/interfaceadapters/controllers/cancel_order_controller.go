package controllers

import "clean-architecture/internal/usecases"

type CancelOrderController struct {
	useCase usecases.CancelOrderInputBoundary
}

func NewCancelOrderController(useCase usecases.CancelOrderInputBoundary) CancelOrderController {
	return CancelOrderController{useCase: useCase}
}

func (c CancelOrderController) Handle(orderID string) error {
	return c.useCase.Execute(usecases.CancelOrderInput{
		OrderID: orderID,
	})
}
