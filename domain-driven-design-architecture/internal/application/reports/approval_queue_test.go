package reports

import (
	"testing"

	"domain-driven-design-architecture/internal/domain/quoting"
)

func TestBuildOrdersAwaitingApprovalReport(t *testing.T) {
	quote, err := quoting.NewQuote("quote-001", "customer-001")
	if err != nil {
		t.Fatal(err)
	}
	price, err := quoting.NewMoney(1000, "USD")
	if err != nil {
		t.Fatal(err)
	}
	line, err := quoting.NewQuoteLineWithCategory("sku-001", quoting.ProductCategoryCustomBuild, 1, price)
	if err != nil {
		t.Fatal(err)
	}
	if err := quote.AddLine(line); err != nil {
		t.Fatal(err)
	}
	if err := quote.SubmitForApproval(quoting.ApprovalDecision{Required: true}); err != nil {
		t.Fatal(err)
	}
	approved, err := quoting.NewQuote("quote-002", "customer-002")
	if err != nil {
		t.Fatal(err)
	}
	if err := approved.AddLine(line); err != nil {
		t.Fatal(err)
	}
	if err := approved.Submit(); err != nil {
		t.Fatal(err)
	}
	report := BuildOrdersAwaitingApprovalReport([]quoting.Quote{quote, approved})
	if len(report.Rows) != 1 || report.Rows[0].QuoteID != "quote-001" || report.Rows[0].TotalCents != 1000 {
		t.Fatalf("unexpected report %+v", report)
	}
}
