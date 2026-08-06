package reporting

import (
	"component-based-architecture/internal/components/orders"
	"component-based-architecture/internal/components/quotes"
	"testing"
)

type quoteReaderStub struct{ items []quotes.QuoteSummary }

func (s quoteReaderStub) GetQuote(quotes.GetQuoteQuery) (quotes.QuoteDetails, error) {
	return quotes.QuoteDetails{}, nil
}
func (s quoteReaderStub) ListQuotes(quotes.ListQuotesQuery) []quotes.QuoteSummary { return s.items }

type orderReaderStub struct{ items []orders.OrderSummary }

func (s orderReaderStub) GetOrder(orders.GetOrderQuery) (orders.OrderDetails, error) {
	return orders.OrderDetails{}, nil
}
func (s orderReaderStub) ListOrders(orders.ListOrdersQuery) []orders.OrderSummary { return s.items }

func TestQuoteConversionReportAggregatesPublishedQueryModels(t *testing.T) {
	component := NewComponent(
		quoteReaderStub{items: []quotes.QuoteSummary{{QuoteID: "quote-001", Status: quotes.QuoteStatusApproved}, {QuoteID: "quote-002", Status: quotes.QuoteStatusDraft}}},
		orderReaderStub{items: []orders.OrderSummary{{QuoteID: "quote-001", Status: orders.OrderStatusPendingPayment}}},
	)
	report := component.QuoteConversionReport()
	if report.TotalQuotes != 2 || report.ApprovedQuotes != 2 || report.ConvertedQuotes != 1 || report.ConversionRate != 0.5 {
		t.Fatalf("unexpected report %+v", report)
	}
}
