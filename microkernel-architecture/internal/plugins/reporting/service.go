package reporting

import "microkernel-architecture/internal/kernel"

type Service struct {
	quotes kernel.QuoteReader
	orders kernel.OrderReader
}

func NewService(quotes kernel.QuoteReader, orders kernel.OrderReader) Service {
	return Service{
		quotes: quotes,
		orders: orders,
	}
}

func (s Service) QuoteConversionReport() (kernel.QuoteConversionReport, error) {
	allQuotes, err := s.quotes.ListQuotes(kernel.ListQuotesQuery{})
	if err != nil {
		return kernel.QuoteConversionReport{}, err
	}

	approvedQuotes, err := s.quotes.ListQuotes(kernel.ListQuotesQuery{Status: "Approved"})
	if err != nil {
		return kernel.QuoteConversionReport{}, err
	}

	convertedOrders, err := s.orders.ListOrders(kernel.ListOrdersQuery{})
	if err != nil {
		return kernel.QuoteConversionReport{}, err
	}

	report := kernel.QuoteConversionReport{
		TotalQuotes:     len(allQuotes),
		ApprovedQuotes:  len(approvedQuotes),
		ConvertedQuotes: len(convertedOrders),
	}

	if report.TotalQuotes > 0 {
		report.ConversionRate = float64(report.ConvertedQuotes) / float64(report.TotalQuotes)
	}

	return report, nil
}
