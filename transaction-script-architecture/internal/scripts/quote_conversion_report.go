package scripts

import "transaction-script-architecture/internal/data"

type QuoteConversionReport struct {
	TotalQuotes     int
	ConvertedQuotes int
	ConversionRate  float64
}

func GetQuoteConversionReport(store *data.Store) (QuoteConversionReport, error) {
	if store == nil {
		return QuoteConversionReport{}, ErrStoreRequired
	}

	report := QuoteConversionReport{TotalQuotes: len(store.Quotes)}
	for _, quote := range store.Quotes {
		if quote.Status == data.QuoteStatusConverted || quote.ConvertedOrderID != "" {
			report.ConvertedQuotes++
		}
	}

	if report.TotalQuotes > 0 {
		report.ConversionRate = float64(report.ConvertedQuotes) / float64(report.TotalQuotes)
	}

	return report, nil
}
