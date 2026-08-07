package reports

import (
	"testing"

	"rich-domain-model-architecture/internal/domain/ordering"
	"rich-domain-model-architecture/internal/domain/quoting"
)

func TestQuoteConversionReportMatchesOrdersToSourceQuotes(t *testing.T) {
	quote := reportQuote(t, "quote-001")
	order, err := ordering.NewOrderFromQuote("order-001", quote)
	if err != nil {
		t.Fatal(err)
	}
	report := BuildQuoteConversionReport([]quoting.Quote{quote}, []ordering.Order{order})
	if report.TotalQuotes != 1 || report.ApprovedQuotes != 1 || report.ConvertedQuotes != 1 || report.ConversionRate != 1 {
		t.Fatalf("report = %+v", report)
	}
}

func TestQuoteConversionReportHandlesNoQuotes(t *testing.T) {
	report := BuildQuoteConversionReport(nil, nil)
	if report.ConversionRate != 0 || report.TotalQuotes != 0 {
		t.Fatalf("report = %+v", report)
	}
}

func reportQuote(t *testing.T, id quoting.QuoteID) quoting.Quote {
	t.Helper()
	price, err := quoting.NewMoney(1000, "USD")
	if err != nil {
		t.Fatal(err)
	}
	line, err := quoting.NewQuoteLine("sku-001", 1, price)
	if err != nil {
		t.Fatal(err)
	}
	quote, err := quoting.NewQuote(id, "customer-001")
	if err != nil {
		t.Fatal(err)
	}
	if err := quote.AddLine(line); err != nil {
		t.Fatal(err)
	}
	if err := quote.SubmitForApproval(quoting.ApprovalDecision{}); err != nil {
		t.Fatal(err)
	}
	return quote
}
