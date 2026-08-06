package reporting

import (
	"component-based-architecture/internal/components/quotes"
	"testing"
)

type approvalQueueQuoteReader struct{}

func (approvalQueueQuoteReader) GetQuote(quotes.GetQuoteQuery) (quotes.QuoteDetails, error) {
	return quotes.QuoteDetails{}, nil
}
func (approvalQueueQuoteReader) ListQuotes(quotes.ListQuotesQuery) []quotes.QuoteSummary {
	return []quotes.QuoteSummary{{QuoteID: "quote-001", CustomerID: "customer-001", Status: quotes.QuoteStatusPendingApproval, LineCount: 2, TotalAmount: 60000}}
}

func TestApprovalQueueIncludesPendingQuoteSummary(t *testing.T) {
	report := NewApprovalQueueComponent(approvalQueueQuoteReader{}).OrdersAwaitingApprovalReport()
	if len(report.Rows) != 1 || report.Rows[0].LineCount != 2 || report.Rows[0].TotalAmount != 60000 {
		t.Fatalf("unexpected report %+v", report)
	}
}
