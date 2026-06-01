package controllers

import "clean-architecture/internal/usecases"

type GetCustomerController struct {
	useCase usecases.GetCustomerInputBoundary
}

func NewGetCustomerController(useCase usecases.GetCustomerInputBoundary) GetCustomerController {
	return GetCustomerController{useCase: useCase}
}

func (c GetCustomerController) Handle(customerID string) error {
	return c.useCase.Execute(usecases.GetCustomerInput{
		CustomerID: customerID,
	})
}
