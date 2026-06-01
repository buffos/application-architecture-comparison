package controllers

import "clean-architecture/internal/usecases"

type QuoteConversionReportController struct {
	useCase usecases.QuoteConversionReportInputBoundary
}

func NewQuoteConversionReportController(useCase usecases.QuoteConversionReportInputBoundary) QuoteConversionReportController {
	return QuoteConversionReportController{useCase: useCase}
}

func (c QuoteConversionReportController) Handle() error {
	return c.useCase.Execute(usecases.QuoteConversionReportInput{})
}
