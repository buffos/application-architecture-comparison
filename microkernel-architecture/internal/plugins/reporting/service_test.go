package reporting

import (
	"testing"

	"microkernel-architecture/internal/kernel"
)

type stubQuoteReader struct {
	list func(query kernel.ListQuotesQuery) ([]kernel.QuoteSummary, error)
}

func (r stubQuoteReader) GetQuote(query kernel.GetQuoteQuery) (kernel.QuoteDetails, error) {
	return kernel.QuoteDetails{}, nil
}

func (r stubQuoteReader) ListQuotes(query kernel.ListQuotesQuery) ([]kernel.QuoteSummary, error) {
	return r.list(query)
}

type stubOrderReader struct {
	list func(query kernel.ListOrdersQuery) ([]kernel.OrderSummary, error)
	get  func(query kernel.GetOrderQuery) (kernel.OrderDetails, error)
}

func (r stubOrderReader) GetOrder(query kernel.GetOrderQuery) (kernel.OrderDetails, error) {
	if r.get != nil {
		return r.get(query)
	}

	return kernel.OrderDetails{}, nil
}

func (r stubOrderReader) ListOrders(query kernel.ListOrdersQuery) ([]kernel.OrderSummary, error) {
	return r.list(query)
}

type stubReturnReader struct {
	list func(query kernel.ListReturnRequestsQuery) ([]kernel.ReturnRequestSummary, error)
	get  func(query kernel.GetReturnRequestQuery) (kernel.ReturnRequestDetails, error)
}

func (r stubReturnReader) GetReturnRequest(query kernel.GetReturnRequestQuery) (kernel.ReturnRequestDetails, error) {
	if r.get != nil {
		return r.get(query)
	}

	return kernel.ReturnRequestDetails{}, nil
}

func (r stubReturnReader) ListReturnRequests(query kernel.ListReturnRequestsQuery) ([]kernel.ReturnRequestSummary, error) {
	if r.list != nil {
		return r.list(query)
	}

	return nil, nil
}

type stubInventoryReader struct {
	list func() ([]kernel.StockSnapshot, error)
}

func (r stubInventoryReader) ListStock() ([]kernel.StockSnapshot, error) {
	if r.list != nil {
		return r.list()
	}

	return nil, nil
}

func TestQuoteConversionReportCombinesQuoteAndOrderCounts(t *testing.T) {
	service := NewService(
		stubQuoteReader{
			list: func(query kernel.ListQuotesQuery) ([]kernel.QuoteSummary, error) {
				if query.Status == "Approved" {
					return []kernel.QuoteSummary{
						{QuoteID: "quote-001", Status: "Approved"},
						{QuoteID: "quote-002", Status: "Approved"},
					}, nil
				}

				return []kernel.QuoteSummary{
					{QuoteID: "quote-001", Status: "Approved"},
					{QuoteID: "quote-002", Status: "Approved"},
					{QuoteID: "quote-003", Status: "Draft"},
				}, nil
			},
		},
		stubOrderReader{
			list: func(query kernel.ListOrdersQuery) ([]kernel.OrderSummary, error) {
				return []kernel.OrderSummary{
					{OrderID: "order-001"},
				}, nil
			},
		},
		stubReturnReader{},
		stubInventoryReader{},
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
		stubQuoteReader{},
		stubOrderReader{
			list: func(query kernel.ListOrdersQuery) ([]kernel.OrderSummary, error) {
				return []kernel.OrderSummary{
					{OrderID: "order-001", Status: "Shipped"},
				}, nil
			},
			get: func(query kernel.GetOrderQuery) (kernel.OrderDetails, error) {
				return kernel.OrderDetails{
					OrderID: "order-001",
					Status:  "Shipped",
					Lines: []kernel.OrderLineDetails{
						{ProductSKU: "sku-001", ProductCategory: "Standard", Quantity: 4},
						{ProductSKU: "sku-002", ProductCategory: "CustomBuild", Quantity: 2},
					},
				}, nil
			},
		},
		stubReturnReader{
			list: func(query kernel.ListReturnRequestsQuery) ([]kernel.ReturnRequestSummary, error) {
				return []kernel.ReturnRequestSummary{
					{ReturnRequestID: "return-001", Status: "Refunded"},
				}, nil
			},
			get: func(query kernel.GetReturnRequestQuery) (kernel.ReturnRequestDetails, error) {
				return kernel.ReturnRequestDetails{
					ReturnRequestID: "return-001",
					Status:          "Refunded",
					Lines: []kernel.ReturnRequestLineDetails{
						{ProductSKU: "sku-001", ProductCategory: "Standard", Quantity: 1},
					},
				}, nil
			},
		},
		stubInventoryReader{},
	)

	report, err := service.ReturnRateByCategoryReport()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(report.Rows) != 2 {
		t.Fatalf("expected two category rows, got %+v", report.Rows)
	}

	if report.Rows[0].Category != "CustomBuild" || report.Rows[0].ShippedQuantity != 2 || report.Rows[0].ReturnedQuantity != 0 || report.Rows[0].ReturnRate != 0 {
		t.Fatalf("unexpected first row %+v", report.Rows[0])
	}

	if report.Rows[1].Category != "Standard" || report.Rows[1].ShippedQuantity != 4 || report.Rows[1].ReturnedQuantity != 1 || report.Rows[1].ReturnRate != 0.25 {
		t.Fatalf("unexpected second row %+v", report.Rows[1])
	}
}

func TestLowStockItemsReportFiltersByThreshold(t *testing.T) {
	service := NewService(
		stubQuoteReader{},
		stubOrderReader{},
		stubReturnReader{},
		stubInventoryReader{
			list: func() ([]kernel.StockSnapshot, error) {
				return []kernel.StockSnapshot{
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

	if report.Rows[0].ProductSKU != "sku-001" || report.Rows[0].Available != 2 {
		t.Fatalf("unexpected first row %+v", report.Rows[0])
	}

	if report.Rows[1].ProductSKU != "sku-003" || report.Rows[1].Available != 5 {
		t.Fatalf("unexpected second row %+v", report.Rows[1])
	}
}

func TestOrdersAwaitingApprovalReportReturnsPendingApprovalQueue(t *testing.T) {
	service := NewService(
		stubQuoteReader{
			list: func(query kernel.ListQuotesQuery) ([]kernel.QuoteSummary, error) {
				if query.Status != "PendingApproval" {
					return nil, nil
				}

				return []kernel.QuoteSummary{
					{QuoteID: "quote-002", CustomerID: "customer-001", Status: "PendingApproval", LineCount: 2, TotalAmount: 60000},
				}, nil
			},
		},
		stubOrderReader{},
		stubReturnReader{},
		stubInventoryReader{},
	)

	report, err := service.OrdersAwaitingApprovalReport()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(report.Rows) != 1 {
		t.Fatalf("expected one pending approval row, got %+v", report.Rows)
	}

	row := report.Rows[0]
	if row.QuoteID != "quote-002" || row.TotalAmount != 60000 || row.LineCount != 2 {
		t.Fatalf("unexpected approval queue row %+v", row)
	}
}
