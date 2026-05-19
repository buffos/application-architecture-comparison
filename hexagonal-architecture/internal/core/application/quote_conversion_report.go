package application

import "hexagonal-architecture/internal/core/ports"

type QuoteConversionReport struct {
	TotalQuotes     int
	ApprovedQuotes  int
	ConvertedQuotes int
	ConversionRate  float64
}

type GetQuoteConversionReportUseCase struct {
	quotes ports.QuoteRepository
	orders ports.OrderRepository
}

func NewGetQuoteConversionReportUseCase(quotes ports.QuoteRepository, orders ports.OrderRepository) GetQuoteConversionReportUseCase {
	return GetQuoteConversionReportUseCase{
		quotes: quotes,
		orders: orders,
	}
}

func (uc GetQuoteConversionReportUseCase) Execute() (QuoteConversionReport, error) {
	quotes, err := uc.quotes.ListByStatus("")
	if err != nil {
		return QuoteConversionReport{}, err
	}

	orders, err := uc.orders.ListByStatus("")
	if err != nil {
		return QuoteConversionReport{}, err
	}

	report := QuoteConversionReport{
		TotalQuotes: len(quotes),
	}

	convertedByQuoteID := make(map[string]struct{}, len(orders))
	for _, order := range orders {
		convertedByQuoteID[order.SourceQuoteID] = struct{}{}
	}

	for _, quote := range quotes {
		if quote.Status == "Approved" {
			report.ApprovedQuotes++
		}
		if _, ok := convertedByQuoteID[quote.ID]; ok {
			report.ConvertedQuotes++
		}
	}

	if report.TotalQuotes > 0 {
		report.ConversionRate = float64(report.ConvertedQuotes) / float64(report.TotalQuotes)
	}

	return report, nil
}
