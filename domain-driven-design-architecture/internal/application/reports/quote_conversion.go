package reports

import (
	"domain-driven-design-architecture/internal/domain/ordering"
	"domain-driven-design-architecture/internal/domain/quoting"
)

type QuoteConversionReport struct {
	TotalQuotes     int
	ApprovedQuotes  int
	ConvertedQuotes int
	ConversionRate  float64
}

func BuildQuoteConversionReport(quotes []quoting.Quote, orders []ordering.Order) QuoteConversionReport {
	report := QuoteConversionReport{TotalQuotes: len(quotes)}
	for _, quote := range quotes {
		if quote.Status() == quoting.QuoteStatusApproved {
			report.ApprovedQuotes++
		}
	}
	converted := make(map[quoting.QuoteID]struct{})
	for _, order := range orders {
		if order.QuoteID() != "" {
			converted[quoting.QuoteID(order.QuoteID())] = struct{}{}
		}
	}
	report.ConvertedQuotes = len(converted)
	if report.TotalQuotes > 0 {
		report.ConversionRate = float64(report.ConvertedQuotes) / float64(report.TotalQuotes)
	}
	return report
}
