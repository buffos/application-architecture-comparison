package controllers

import "clean-architecture/internal/usecases"

type GetProductController struct {
	useCase usecases.GetProductInputBoundary
}

func NewGetProductController(useCase usecases.GetProductInputBoundary) GetProductController {
	return GetProductController{useCase: useCase}
}

func (c GetProductController) Handle(sku string) error {
	return c.useCase.Execute(usecases.GetProductInput{SKU: sku})
}
