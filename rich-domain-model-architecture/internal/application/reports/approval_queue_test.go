package reports

import (
	"testing"

	"rich-domain-model-architecture/internal/domain/quoting"
)

func TestApprovalQueueReportSelectsPendingQuotes(t *testing.T) {
	pending := pendingQuoteForReport(t)
	approved := reportQuote(t, "quote-approved")
	report := BuildOrdersAwaitingApprovalReport([]quoting.Quote{approved, pending})
	if len(report.Rows) != 1 {
		t.Fatalf("rows = %+v", report.Rows)
	}
	row := report.Rows[0]
	if row.QuoteID != string(pending.ID()) || row.TotalCents != 45000 || row.LineCount != 1 {
		t.Fatalf("row = %+v", row)
	}
	if pending.Status() != quoting.QuoteStatusPendingApproval {
		t.Fatal("report changed pending quote")
	}
}

func pendingQuoteForReport(t *testing.T) quoting.Quote {
	t.Helper()
	price, err := quoting.NewMoney(45000, "USD")
	if err != nil {
		t.Fatal(err)
	}
	line, err := quoting.NewQuoteLineWithCategory("sku-custom", quoting.ProductCategoryCustomBuild, 1, price)
	if err != nil {
		t.Fatal(err)
	}
	quote, err := quoting.NewQuote("quote-pending", "customer-001")
	if err != nil {
		t.Fatal(err)
	}
	if err := quote.AddLine(line); err != nil {
		t.Fatal(err)
	}
	if err := quote.SubmitForApproval(quoting.ApprovalDecision{Required: true}); err != nil {
		t.Fatal(err)
	}
	return quote
}
