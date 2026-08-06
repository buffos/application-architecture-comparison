package quoting

import "testing"

func TestQuoteApprovalServiceReturnsCustomBuildReason(t *testing.T) {
	quote, err := NewQuote("quote-001", "customer-001")
	if err != nil {
		t.Fatal(err)
	}
	price, err := NewMoney(45000, "USD")
	if err != nil {
		t.Fatal(err)
	}
	line, err := NewQuoteLineWithCategory("sku-custom", ProductCategoryCustomBuild, 1, price)
	if err != nil {
		t.Fatal(err)
	}
	if err := quote.AddLine(line); err != nil {
		t.Fatal(err)
	}
	decision := NewQuoteApprovalService().Evaluate(quote)
	if !decision.Required || len(decision.Reasons) != 1 || decision.Reasons[0].Code != ApprovalReasonCustomBuild {
		t.Fatalf("unexpected decision %+v", decision)
	}
	if quote.Status() != QuoteStatusDraft || len(quote.Lines()) != 1 {
		t.Fatalf("approval evaluation mutated quote: status=%s lines=%d", quote.Status(), len(quote.Lines()))
	}
}

func TestQuoteApprovalServiceAllowsStandardQuote(t *testing.T) {
	quote, err := NewQuote("quote-001", "customer-001")
	if err != nil {
		t.Fatal(err)
	}
	price, err := NewMoney(15000, "USD")
	if err != nil {
		t.Fatal(err)
	}
	line, err := NewQuoteLine("sku-standard", 1, price)
	if err != nil {
		t.Fatal(err)
	}
	if err := quote.AddLine(line); err != nil {
		t.Fatal(err)
	}
	decision := NewQuoteApprovalService().Evaluate(quote)
	if decision.Required || len(decision.Reasons) != 0 {
		t.Fatalf("unexpected decision %+v", decision)
	}
}
