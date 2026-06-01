package controllers

import "clean-architecture/internal/usecases"

type RejectReturnController struct {
	useCase usecases.RejectReturnInputBoundary
}

func NewRejectReturnController(useCase usecases.RejectReturnInputBoundary) RejectReturnController {
	return RejectReturnController{useCase: useCase}
}

func (c RejectReturnController) Handle(returnRequestID string, idempotencyKey string, reviewedBy string, reviewNote string) error {
	return c.useCase.Execute(usecases.RejectReturnInput{
		ReturnRequestID: returnRequestID,
		IdempotencyKey:  idempotencyKey,
		ReviewedBy:      reviewedBy,
		ReviewNote:      reviewNote,
	})
}
