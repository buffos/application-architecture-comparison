package reporting

import (
	"testing"

	"modular-monolith/internal/modules/orders"
	"modular-monolith/internal/modules/quotes"
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
