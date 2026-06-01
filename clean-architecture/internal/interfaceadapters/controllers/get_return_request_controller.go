package controllers

import "clean-architecture/internal/usecases"

type GetReturnRequestController struct {
	useCase usecases.GetReturnRequestInputBoundary
}

func NewGetReturnRequestController(useCase usecases.GetReturnRequestInputBoundary) GetReturnRequestController {
	return GetReturnRequestController{useCase: useCase}
}

func (c GetReturnRequestController) Handle(returnRequestID string) error {
	return c.useCase.Execute(usecases.GetReturnRequestInput{
		ReturnRequestID: returnRequestID,
	})
}
