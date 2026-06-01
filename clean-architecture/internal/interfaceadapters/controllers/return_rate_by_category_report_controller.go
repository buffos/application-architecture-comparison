package controllers

import "clean-architecture/internal/usecases"

type ReturnRateByCategoryReportController struct {
	useCase usecases.ReturnRateByCategoryReportInputBoundary
}

func NewReturnRateByCategoryReportController(useCase usecases.ReturnRateByCategoryReportInputBoundary) ReturnRateByCategoryReportController {
	return ReturnRateByCategoryReportController{useCase: useCase}
}

func (c ReturnRateByCategoryReportController) Handle() error {
	return c.useCase.Execute(usecases.ReturnRateByCategoryReportInput{})
}
