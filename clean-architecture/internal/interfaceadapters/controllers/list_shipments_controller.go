package controllers

import "clean-architecture/internal/usecases"

type ListShipmentsController struct {
	useCase usecases.ListShipmentsInputBoundary
}

func NewListShipmentsController(useCase usecases.ListShipmentsInputBoundary) ListShipmentsController {
	return ListShipmentsController{useCase: useCase}
}

func (c ListShipmentsController) Handle(orderID string) error {
	return c.useCase.Execute(usecases.ListShipmentsInput{
		OrderID: orderID,
	})
}
