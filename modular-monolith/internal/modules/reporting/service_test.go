package reporting

import (
	"testing"

	"modular-monolith/internal/modules/inventory"
	"modular-monolith/internal/modules/orders"
	"modular-monolith/internal/modules/quotes"
	"modular-monolith/internal/modules/returns"
)

type stubQuoteReader struct {
	list func(query quotes.ListQuotesQuery) ([]quotes.QuoteDetails, error)
}

func (r stubQuoteReader) ListQuotes(query quotes.ListQuotesQuery) ([]quotes.QuoteDetails, error) {
	return r.list(query)
}

type stubOrderReader struct {
	list func(query orders.ListOrdersQuery) ([]orders.OrderDetails, error)
}

func (r stubOrderReader) ListOrders(query orders.ListOrdersQuery) ([]orders.OrderDetails, error) {
	return r.list(query)
}

type stubReturnReader struct {
	list func(query returns.ListReturnRequestsQuery) ([]returns.ReturnRequestDetails, error)
}

func (r stubReturnReader) ListReturnRequests(query returns.ListReturnRequestsQuery) ([]returns.ReturnRequestDetails, error) {
	return r.list(query)
}

type stubInventoryReader struct {
	list func() ([]inventory.StockSnapshot, error)
}

func (r stubInventoryReader) ListStock() ([]inventory.StockSnapshot, error) {
	return r.list()
}

func TestQuoteConversionReportCombinesQuoteAndOrderCounts(t *testing.T) {
	service := NewService(
		stubQuoteReader{
			list: func(query quotes.ListQuotesQuery) ([]quotes.QuoteDetails, error) {
				if query.Status == quotes.QuoteStatusApproved {
					return []quotes.QuoteDetails{
						{QuoteID: "quote-001", Status: quotes.QuoteStatusApproved},
						{QuoteID: "quote-002", Status: quotes.QuoteStatusApproved},
					}, nil
				}

				return []quotes.QuoteDetails{
					{QuoteID: "quote-001", Status: quotes.QuoteStatusApproved},
					{QuoteID: "quote-002", Status: quotes.QuoteStatusApproved},
					{QuoteID: "quote-003", Status: quotes.QuoteStatusDraft},
				}, nil
			},
		},
		stubOrderReader{
			list: func(query orders.ListOrdersQuery) ([]orders.OrderDetails, error) {
				return []orders.OrderDetails{
					{OrderID: "order-001"},
				}, nil
			},
		},
		stubReturnReader{
			list: func(query returns.ListReturnRequestsQuery) ([]returns.ReturnRequestDetails, error) {
				return nil, nil
			},
		},
		stubInventoryReader{
			list: func() ([]inventory.StockSnapshot, error) {
				return nil, nil
			},
		},
	)

	report, err := service.QuoteConversionReport()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if report.TotalQuotes != 3 {
		t.Fatalf("expected total quotes 3, got %d", report.TotalQuotes)
	}

	if report.ApprovedQuotes != 2 {
		t.Fatalf("expected approved quotes 2, got %d", report.ApprovedQuotes)
	}

	if report.ConvertedQuotes != 1 {
		t.Fatalf("expected converted quotes 1, got %d", report.ConvertedQuotes)
	}

	if report.ConversionRate != 1.0/3.0 {
		t.Fatalf("expected conversion rate 1/3, got %f", report.ConversionRate)
	}
}

func TestReturnRateByCategoryReportGroupsShippedAndReturnedQuantities(t *testing.T) {
	service := NewService(
		stubQuoteReader{list: func(query quotes.ListQuotesQuery) ([]quotes.QuoteDetails, error) { return nil, nil }},
		stubOrderReader{
			list: func(query orders.ListOrdersQuery) ([]orders.OrderDetails, error) {
				return []orders.OrderDetails{
					{
						OrderID: "order-001",
						Status:  orders.OrderStatusShipped,
						Lines: []orders.OrderLineDetails{
							{ProductSKU: "sku-001", ProductCategory: "Standard", Quantity: 4},
							{ProductSKU: "sku-002", ProductCategory: "CustomBuild", Quantity: 2},
						},
					},
				}, nil
			},
		},
		stubReturnReader{
			list: func(query returns.ListReturnRequestsQuery) ([]returns.ReturnRequestDetails, error) {
				return []returns.ReturnRequestDetails{
					{
						ReturnRequestID: "return-001",
						Status:          returns.ReturnRequestStatusRefunded,
						Lines: []returns.ReturnLineDetails{
							{ProductSKU: "sku-001", ProductCategory: "Standard", Quantity: 1},
						},
					},
				}, nil
			},
		},
		stubInventoryReader{
			list: func() ([]inventory.StockSnapshot, error) {
				return nil, nil
			},
		},
	)

	report, err := service.ReturnRateByCategoryReport()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(report.Rows) != 2 {
		t.Fatalf("expected two category rows, got %+v", report.Rows)
	}

	var standard *ReturnRateByCategoryRow
	for i := range report.Rows {
		if report.Rows[i].Category == "Standard" {
			standard = &report.Rows[i]
		}
	}

	if standard == nil {
		t.Fatalf("expected Standard category row, got %+v", report.Rows)
	}

	if standard.ShippedQuantity != 4 || standard.ReturnedQuantity != 1 {
		t.Fatalf("expected Standard row 4 shipped / 1 returned, got %+v", *standard)
	}

	if standard.ReturnRate != 0.25 {
		t.Fatalf("expected Standard return rate 0.25, got %f", standard.ReturnRate)
	}
}

func TestLowStockItemsReportFiltersByThreshold(t *testing.T) {
	service := NewService(
		stubQuoteReader{list: func(query quotes.ListQuotesQuery) ([]quotes.QuoteDetails, error) { return nil, nil }},
		stubOrderReader{list: func(query orders.ListOrdersQuery) ([]orders.OrderDetails, error) { return nil, nil }},
		stubReturnReader{list: func(query returns.ListReturnRequestsQuery) ([]returns.ReturnRequestDetails, error) { return nil, nil }},
		stubInventoryReader{
			list: func() ([]inventory.StockSnapshot, error) {
				return []inventory.StockSnapshot{
					{ProductSKU: "sku-001", Available: 2},
					{ProductSKU: "sku-002", Available: 7},
					{ProductSKU: "sku-003", Available: 5},
				}, nil
			},
		},
	)

	report, err := service.LowStockItemsReport(5)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(report.Rows) != 2 {
		t.Fatalf("expected two low stock rows, got %+v", report.Rows)
	}
}
