package reporting

import (
	"modular-monolith/internal/modules/orders"
	"modular-monolith/internal/modules/quotes"
)

type QuoteReader interface {
	ListQuotes(query quotes.ListQuotesQuery) ([]quotes.QuoteDetails, error)
}

type OrderReader interface {
	ListOrders(query orders.ListOrdersQuery) ([]orders.OrderDetails, error)
}

type Service struct {
	quotes QuoteReader
	orders OrderReader
}

func NewService(quotes QuoteReader, orders OrderReader) Service {
	return Service{
		quotes: quotes,
		orders: orders,
	}
}

type QuoteConversionReport struct {
	TotalQuotes     int
	ApprovedQuotes  int
	ConvertedQuotes int
	ConversionRate  float64
}

func (s Service) QuoteConversionReport() (QuoteConversionReport, error) {
	allQuotes, err := s.quotes.ListQuotes(quotes.ListQuotesQuery{})
	if err != nil {
		return QuoteConversionReport{}, err
	}

	approvedQuotes, err := s.quotes.ListQuotes(quotes.ListQuotesQuery{Status: quotes.QuoteStatusApproved})
	if err != nil {
		return QuoteConversionReport{}, err
	}

	convertedOrders, err := s.orders.ListOrders(orders.ListOrdersQuery{})
	if err != nil {
		return QuoteConversionReport{}, err
	}

	report := QuoteConversionReport{
		TotalQuotes:     len(allQuotes),
		ApprovedQuotes:  len(approvedQuotes),
		ConvertedQuotes: len(convertedOrders),
	}

	if report.TotalQuotes > 0 {
		report.ConversionRate = float64(report.ConvertedQuotes) / float64(report.TotalQuotes)
	}

	return report, nil
}
