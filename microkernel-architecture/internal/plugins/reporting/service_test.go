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
}

func (r stubOrderReader) GetOrder(query kernel.GetOrderQuery) (kernel.OrderDetails, error) {
	return kernel.OrderDetails{}, nil
}

func (r stubOrderReader) ListOrders(query kernel.ListOrdersQuery) ([]kernel.OrderSummary, error) {
	return r.list(query)
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
