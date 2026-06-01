package presenters

import (
	"fmt"

	"clean-architecture/internal/usecases"
)

type QuoteConversionReportViewModel struct {
	Message         string
	TotalQuotes     int
	ApprovedQuotes  int
	ConvertedQuotes int
	ConversionRate  float64
}

type QuoteConversionReportPresenter struct {
	viewModel QuoteConversionReportViewModel
}

func NewQuoteConversionReportPresenter() *QuoteConversionReportPresenter {
	return &QuoteConversionReportPresenter{}
}

func (p *QuoteConversionReportPresenter) Present(output usecases.QuoteConversionReportOutput) error {
	p.viewModel = QuoteConversionReportViewModel{
		Message:         fmt.Sprintf("quote conversion report: total=%d approved=%d converted=%d rate=%.2f", output.TotalQuotes, output.ApprovedQuotes, output.ConvertedQuotes, output.ConversionRate),
		TotalQuotes:     output.TotalQuotes,
		ApprovedQuotes:  output.ApprovedQuotes,
		ConvertedQuotes: output.ConvertedQuotes,
		ConversionRate:  output.ConversionRate,
	}

	return nil
}

func (p *QuoteConversionReportPresenter) ViewModel() QuoteConversionReportViewModel {
	return p.viewModel
}
