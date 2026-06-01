package controllers

import "clean-architecture/internal/usecases"

type ApprovePaymentReviewController struct {
	useCase usecases.ApprovePaymentReviewInputBoundary
}

func NewApprovePaymentReviewController(useCase usecases.ApprovePaymentReviewInputBoundary) ApprovePaymentReviewController {
	return ApprovePaymentReviewController{useCase: useCase}
}

func (c ApprovePaymentReviewController) Handle(orderID string) error {
	return c.useCase.Execute(usecases.ApprovePaymentReviewInput{
		OrderID: orderID,
	})
}
