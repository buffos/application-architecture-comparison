package reports

import (
	"testing"

	"domain-driven-design-architecture/internal/domain/ordering"
	"domain-driven-design-architecture/internal/domain/quoting"
)

func TestBuildQuoteConversionReport(t *testing.T) {
	quote, err := quoting.NewQuote("quote-001", "customer-001")
	if err != nil {
		t.Fatal(err)
	}
	price, err := quoting.NewMoney(1000, "USD")
	if err != nil {
		t.Fatal(err)
	}
	line, err := quoting.NewQuoteLine("sku-001", 1, price)
	if err != nil {
		t.Fatal(err)
	}
	if err := quote.AddLine(line); err != nil {
		t.Fatal(err)
	}
	if err := quote.Submit(); err != nil {
		t.Fatal(err)
	}
	order, err := ordering.NewOrderFromQuote("order-001", quote)
	if err != nil {
		t.Fatal(err)
	}
	report := BuildQuoteConversionReport([]quoting.Quote{quote}, []ordering.Order{order, order})
	if report.TotalQuotes != 1 || report.ApprovedQuotes != 1 || report.ConvertedQuotes != 1 || report.ConversionRate != 1 {
		t.Fatalf("unexpected report %+v", report)
	}
}
