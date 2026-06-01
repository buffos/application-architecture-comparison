package controllers

import "clean-architecture/internal/usecases"

type OrdersAwaitingApprovalReportController struct {
	useCase usecases.OrdersAwaitingApprovalReportInputBoundary
}

func NewOrdersAwaitingApprovalReportController(useCase usecases.OrdersAwaitingApprovalReportInputBoundary) OrdersAwaitingApprovalReportController {
	return OrdersAwaitingApprovalReportController{useCase: useCase}
}

func (c OrdersAwaitingApprovalReportController) Handle() error {
	return c.useCase.Execute(usecases.OrdersAwaitingApprovalReportInput{})
}
