package controllers

import "clean-architecture/internal/usecases"

type CapturePaymentController struct {
	useCase usecases.CapturePaymentInputBoundary
}

func NewCapturePaymentController(useCase usecases.CapturePaymentInputBoundary) CapturePaymentController {
	return CapturePaymentController{useCase: useCase}
}

func (c CapturePaymentController) Handle(orderID string) error {
	return c.useCase.Execute(usecases.CapturePaymentInput{
		OrderID: orderID,
	})
}
