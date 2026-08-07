package records

// QuoteConversionReport is a read-time projection of quote lifecycle data.
type QuoteConversionReport struct {
	TotalQuotes     int
	ConvertedQuotes int
	ConversionRate  float64
}

// GetQuoteConversionReport scans current quote rows without changing them.
func GetQuoteConversionReport(db *Database) (QuoteConversionReport, error) {
	if db == nil {
		return QuoteConversionReport{}, ErrDatabaseRequired
	}

	report := QuoteConversionReport{TotalQuotes: len(db.quotes)}
	for _, row := range db.quotes {
		if row.Status == QuoteStatusConverted || row.ConvertedOrderID != "" {
			report.ConvertedQuotes++
		}
	}
	if report.TotalQuotes > 0 {
		report.ConversionRate = float64(report.ConvertedQuotes) / float64(report.TotalQuotes)
	}
	return report, nil
}
