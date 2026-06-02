package reporting

import (
	"modular-monolith/internal/modules/inventory"
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

type InventoryReader interface {
	ListStock() ([]inventory.StockSnapshot, error)
}

type Service struct {
	quotes    QuoteReader
	orders    OrderReader
	returns   ReturnReader
	inventory InventoryReader
}

func NewService(quotes QuoteReader, orders OrderReader, returns ReturnReader, inventory InventoryReader) Service {
	return Service{
		quotes:    quotes,
		orders:    orders,
		returns:   returns,
		inventory: inventory,
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

type LowStockItemsReportRow struct {
	ProductSKU string
	Available  int
}

type LowStockItemsReport struct {
	Rows []LowStockItemsReportRow
}

type OrdersAwaitingApprovalRow struct {
	QuoteID     string
	CustomerID  string
	LineCount   int
	TotalAmount int
}

type OrdersAwaitingApprovalReport struct {
	Rows []OrdersAwaitingApprovalRow
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

func (s Service) LowStockItemsReport(threshold int) (LowStockItemsReport, error) {
	stock, err := s.inventory.ListStock()
	if err != nil {
		return LowStockItemsReport{}, err
	}

	report := LowStockItemsReport{
		Rows: make([]LowStockItemsReportRow, 0),
	}
	for _, item := range stock {
		if item.Available <= threshold {
			report.Rows = append(report.Rows, LowStockItemsReportRow{
				ProductSKU: item.ProductSKU,
				Available:  item.Available,
			})
		}
	}

	return report, nil
}

func (s Service) OrdersAwaitingApprovalReport() (OrdersAwaitingApprovalReport, error) {
	pendingQuotes, err := s.quotes.ListQuotes(quotes.ListQuotesQuery{Status: quotes.QuoteStatusPendingApproval})
	if err != nil {
		return OrdersAwaitingApprovalReport{}, err
	}

	report := OrdersAwaitingApprovalReport{
		Rows: make([]OrdersAwaitingApprovalRow, 0, len(pendingQuotes)),
	}
	for _, quote := range pendingQuotes {
		report.Rows = append(report.Rows, OrdersAwaitingApprovalRow{
			QuoteID:     quote.QuoteID,
			CustomerID:  quote.CustomerID,
			LineCount:   quote.LineCount,
			TotalAmount: quote.TotalAmount,
		})
	}

	return report, nil
}
