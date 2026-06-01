package controllers

import "clean-architecture/internal/usecases"

type CreateShipmentController struct {
	useCase usecases.CreateShipmentInputBoundary
}

func NewCreateShipmentController(useCase usecases.CreateShipmentInputBoundary) CreateShipmentController {
	return CreateShipmentController{useCase: useCase}
}

func (c CreateShipmentController) Handle(orderID string) error {
	return c.useCase.Execute(usecases.CreateShipmentInput{
		OrderID: orderID,
	})
}
