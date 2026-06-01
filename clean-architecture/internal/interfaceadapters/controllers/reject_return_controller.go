package controllers

import "clean-architecture/internal/usecases"

type RejectReturnController struct {
	useCase usecases.RejectReturnInputBoundary
}

func NewRejectReturnController(useCase usecases.RejectReturnInputBoundary) RejectReturnController {
	return RejectReturnController{useCase: useCase}
}

func (c RejectReturnController) Handle(returnRequestID string, reviewedBy string, reviewNote string) error {
	return c.useCase.Execute(usecases.RejectReturnInput{
		ReturnRequestID: returnRequestID,
		ReviewedBy:      reviewedBy,
		ReviewNote:      reviewNote,
	})
}
