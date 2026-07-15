package reporting

import (
	"slices"

	"microkernel-architecture/internal/kernel"
	"microkernel-architecture/internal/plugins/orders"
	"microkernel-architecture/internal/plugins/returns"
)

type Service struct {
	quotes  kernel.QuoteReader
	orders  kernel.OrderReader
	returns kernel.ReturnReader
}

func NewService(quotes kernel.QuoteReader, orders kernel.OrderReader, returns kernel.ReturnReader) Service {
	return Service{
		quotes:  quotes,
		orders:  orders,
		returns: returns,
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

func (s Service) ReturnRateByCategoryReport() (kernel.ReturnRateByCategoryReport, error) {
	shippedOrders, err := s.orders.ListOrders(kernel.ListOrdersQuery{Status: orders.OrderStatusShipped})
	if err != nil {
		return kernel.ReturnRateByCategoryReport{}, err
	}

	refundedReturns, err := s.returns.ListReturnRequests(kernel.ListReturnRequestsQuery{Status: returns.ReturnRequestStatusRefunded})
	if err != nil {
		return kernel.ReturnRateByCategoryReport{}, err
	}

	rowsByCategory := map[string]*kernel.ReturnRateByCategoryRow{}

	for _, order := range shippedOrders {
		details, err := s.orders.GetOrder(kernel.GetOrderQuery{OrderID: order.OrderID})
		if err != nil {
			return kernel.ReturnRateByCategoryReport{}, err
		}

		for _, line := range details.Lines {
			row := ensureRow(rowsByCategory, line.ProductCategory)
			row.ShippedQuantity += line.Quantity
		}
	}

	for _, request := range refundedReturns {
		details, err := s.returns.GetReturnRequest(kernel.GetReturnRequestQuery{ReturnRequestID: request.ReturnRequestID})
		if err != nil {
			return kernel.ReturnRateByCategoryReport{}, err
		}

		for _, line := range details.Lines {
			row := ensureRow(rowsByCategory, line.ProductCategory)
			row.ReturnedQuantity += line.Quantity
		}
	}

	report := kernel.ReturnRateByCategoryReport{
		Rows: make([]kernel.ReturnRateByCategoryRow, 0, len(rowsByCategory)),
	}
	for _, row := range rowsByCategory {
		if row.ShippedQuantity > 0 {
			row.ReturnRate = float64(row.ReturnedQuantity) / float64(row.ShippedQuantity)
		}
		report.Rows = append(report.Rows, *row)
	}

	slices.SortFunc(report.Rows, func(a, b kernel.ReturnRateByCategoryRow) int {
		if a.Category < b.Category {
			return -1
		}
		if a.Category > b.Category {
			return 1
		}
		return 0
	})

	return report, nil
}

func ensureRow(rowsByCategory map[string]*kernel.ReturnRateByCategoryRow, category string) *kernel.ReturnRateByCategoryRow {
	row := rowsByCategory[category]
	if row == nil {
		row = &kernel.ReturnRateByCategoryRow{Category: category}
		rowsByCategory[category] = row
	}

	return row
}
