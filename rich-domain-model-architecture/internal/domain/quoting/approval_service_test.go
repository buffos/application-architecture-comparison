package quoting

import "testing"

func TestQuoteApprovalServiceEvaluatesPolicyWithoutMutation(t *testing.T) {
	price, err := NewMoney(45000, "USD")
	if err != nil {
		t.Fatal(err)
	}
	line, err := NewQuoteLineFromProductSnapshotWithCategory("sku-custom", "Custom Desk", ProductCategoryCustomBuild, 1, price)
	if err != nil {
		t.Fatal(err)
	}
	quote, err := NewQuote("quote-001", "customer-001")
	if err != nil {
		t.Fatal(err)
	}
	if err := quote.AddLine(line); err != nil {
		t.Fatal(err)
	}

	decision := NewQuoteApprovalService().Evaluate(quote)
	if !decision.Required || len(decision.Reasons) != 1 {
		t.Fatalf("decision = %+v, want one required reason", decision)
	}
	if decision.Reasons[0].Code != ApprovalReasonCustomBuild {
		t.Fatalf("reason code = %s, want %s", decision.Reasons[0].Code, ApprovalReasonCustomBuild)
	}
	if quote.Status() != QuoteStatusDraft || len(quote.Lines()) != 1 {
		t.Fatal("approval evaluation mutated the quote")
	}
}

func TestQuoteApprovalServiceDoesNotRequireApprovalForStandardLines(t *testing.T) {
	price, err := NewMoney(1000, "USD")
	if err != nil {
		t.Fatal(err)
	}
	line, err := NewQuoteLine("sku-standard", 1, price)
	if err != nil {
		t.Fatal(err)
	}
	quote, err := NewQuote("quote-001", "customer-001")
	if err != nil {
		t.Fatal(err)
	}
	if err := quote.AddLine(line); err != nil {
		t.Fatal(err)
	}

	decision := NewQuoteApprovalService().Evaluate(quote)
	if decision.Required || len(decision.Reasons) != 0 {
		t.Fatalf("decision = %+v, want no approval", decision)
	}
}
