package controllers

import "clean-architecture/internal/usecases"

type GetShipmentController struct {
	useCase usecases.GetShipmentInputBoundary
}

func NewGetShipmentController(useCase usecases.GetShipmentInputBoundary) GetShipmentController {
	return GetShipmentController{useCase: useCase}
}

func (c GetShipmentController) Handle(shipmentID string) error {
	return c.useCase.Execute(usecases.GetShipmentInput{
		ShipmentID: shipmentID,
	})
}
