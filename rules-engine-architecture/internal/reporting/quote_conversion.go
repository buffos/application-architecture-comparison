package reporting

import (
	"strings"

	"rules-engine-architecture/internal/engine"
)

type QuoteConversionReport struct {
	TotalQuotes           int
	ConvertedQuotes       int
	ConversionRatePercent float64
}

func BuildQuoteConversionReport(quotes []engine.QuoteFact) QuoteConversionReport {
	report := QuoteConversionReport{TotalQuotes: len(quotes)}
	for _, quote := range quotes {
		if strings.EqualFold(quote.Status, "converted") {
			report.ConvertedQuotes++
		}
	}

	if report.TotalQuotes > 0 {
		report.ConversionRatePercent = float64(report.ConvertedQuotes) * 100 / float64(report.TotalQuotes)
	}
	return report
}
