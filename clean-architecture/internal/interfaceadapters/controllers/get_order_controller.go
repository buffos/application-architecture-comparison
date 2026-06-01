package controllers

import "clean-architecture/internal/usecases"

type GetOrderController struct {
	useCase usecases.GetOrderInputBoundary
}

func NewGetOrderController(useCase usecases.GetOrderInputBoundary) GetOrderController {
	return GetOrderController{useCase: useCase}
}

func (c GetOrderController) Handle(orderID string) error {
	return c.useCase.Execute(usecases.GetOrderInput{
		OrderID: orderID,
	})
}
