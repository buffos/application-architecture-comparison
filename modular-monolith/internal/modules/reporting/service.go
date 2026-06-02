package reporting

import (
	"modular-monolith/internal/modules/orders"
	"modular-monolith/internal/modules/quotes"
	"modular-monolith/internal/modules/returns"
)

type QuoteReader interface {
	ListQuotes(query quotes.ListQuotesQuery) ([]quotes.QuoteDetails, error)
}

type OrderReader interface {
	ListOrders(query orders.ListOrdersQuery) ([]orders.OrderDetails, error)
}

type ReturnReader interface {
	ListReturnRequests(query returns.ListReturnRequestsQuery) ([]returns.ReturnRequestDetails, error)
}

type Service struct {
	quotes  QuoteReader
	orders  OrderReader
	returns ReturnReader
}

func NewService(quotes QuoteReader, orders OrderReader, returns ReturnReader) Service {
	return Service{
		quotes:  quotes,
		orders:  orders,
		returns: returns,
	}
}

type QuoteConversionReport struct {
	TotalQuotes     int
	ApprovedQuotes  int
	ConvertedQuotes int
	ConversionRate  float64
}

type ReturnRateByCategoryRow struct {
	Category         string
	ShippedQuantity  int
	ReturnedQuantity int
	ReturnRate       float64
}

type ReturnRateByCategoryReport struct {
	Rows []ReturnRateByCategoryRow
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

func (s Service) ReturnRateByCategoryReport() (ReturnRateByCategoryReport, error) {
	shippedOrders, err := s.orders.ListOrders(orders.ListOrdersQuery{Status: orders.OrderStatusShipped})
	if err != nil {
		return ReturnRateByCategoryReport{}, err
	}

	refundedReturns, err := s.returns.ListReturnRequests(returns.ListReturnRequestsQuery{Status: returns.ReturnRequestStatusRefunded})
	if err != nil {
		return ReturnRateByCategoryReport{}, err
	}

	rowsByCategory := map[string]*ReturnRateByCategoryRow{}

	for _, order := range shippedOrders {
		for _, line := range order.Lines {
			row := rowsByCategory[line.ProductCategory]
			if row == nil {
				row = &ReturnRateByCategoryRow{Category: line.ProductCategory}
				rowsByCategory[line.ProductCategory] = row
			}
			row.ShippedQuantity += line.Quantity
		}
	}

	for _, request := range refundedReturns {
		for _, line := range request.Lines {
			row := rowsByCategory[line.ProductCategory]
			if row == nil {
				row = &ReturnRateByCategoryRow{Category: line.ProductCategory}
				rowsByCategory[line.ProductCategory] = row
			}
			row.ReturnedQuantity += line.Quantity
		}
	}

	report := ReturnRateByCategoryReport{
		Rows: make([]ReturnRateByCategoryRow, 0, len(rowsByCategory)),
	}
	for _, row := range rowsByCategory {
		if row.ShippedQuantity > 0 {
			row.ReturnRate = float64(row.ReturnedQuantity) / float64(row.ShippedQuantity)
		}
		report.Rows = append(report.Rows, *row)
	}

	return report, nil
}
