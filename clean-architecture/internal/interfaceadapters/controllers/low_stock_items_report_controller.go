package controllers

import "clean-architecture/internal/usecases"

type LowStockItemsReportController struct {
	useCase usecases.LowStockItemsReportInputBoundary
}

func NewLowStockItemsReportController(useCase usecases.LowStockItemsReportInputBoundary) LowStockItemsReportController {
	return LowStockItemsReportController{useCase: useCase}
}

func (c LowStockItemsReportController) Handle(threshold int) error {
	return c.useCase.Execute(usecases.LowStockItemsReportInput{
		Threshold: threshold,
	})
}
