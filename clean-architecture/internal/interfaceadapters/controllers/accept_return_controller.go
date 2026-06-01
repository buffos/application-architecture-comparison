package controllers

import "clean-architecture/internal/usecases"

type AcceptReturnController struct {
	useCase usecases.AcceptReturnInputBoundary
}

func NewAcceptReturnController(useCase usecases.AcceptReturnInputBoundary) AcceptReturnController {
	return AcceptReturnController{useCase: useCase}
}

func (c AcceptReturnController) Handle(returnRequestID string, reviewedBy string, processedBy string) error {
	return c.useCase.Execute(usecases.AcceptReturnInput{
		ReturnRequestID: returnRequestID,
		ReviewedBy:      reviewedBy,
		ProcessedBy:     processedBy,
	})
}
