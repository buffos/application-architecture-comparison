package controllers

import "clean-architecture/internal/usecases"

type RequestReturnController struct {
	useCase usecases.RequestReturnInputBoundary
}

func NewRequestReturnController(useCase usecases.RequestReturnInputBoundary) RequestReturnController {
	return RequestReturnController{useCase: useCase}
}

func (c RequestReturnController) Handle(orderID string, reason string, requestedBy string) error {
	return c.useCase.Execute(usecases.RequestReturnInput{
		OrderID:     orderID,
		Reason:      reason,
		RequestedBy: requestedBy,
	})
}

func (c RequestReturnController) HandlePartial(orderID string, reason string, lines []usecases.RequestReturnLineInput, requestedBy string) error {
	return c.useCase.Execute(usecases.RequestReturnInput{
		OrderID:     orderID,
		Reason:      reason,
		Lines:       lines,
		RequestedBy: requestedBy,
	})
}
