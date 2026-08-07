package quoting

import (
	"errors"
	"testing"
)

func TestQuoteAppliesApprovalDecisionToItsLifecycle(t *testing.T) {
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
	if err := quote.SubmitForApproval(ApprovalDecision{}); err != nil {
		t.Fatal(err)
	}
	if quote.Status() != QuoteStatusApproved {
		t.Fatalf("status = %s, want %s", quote.Status(), QuoteStatusApproved)
	}

	pendingQuote, err := NewQuote("quote-002", "customer-001")
	if err != nil {
		t.Fatal(err)
	}
	customLine, err := NewQuoteLineWithCategory("sku-custom", ProductCategoryCustomBuild, 1, price)
	if err != nil {
		t.Fatal(err)
	}
	if err := pendingQuote.AddLine(customLine); err != nil {
		t.Fatal(err)
	}
	if err := pendingQuote.SubmitForApproval(ApprovalDecision{Required: true}); err != nil {
		t.Fatal(err)
	}
	if pendingQuote.Status() != QuoteStatusPendingApproval {
		t.Fatalf("status = %s, want %s", pendingQuote.Status(), QuoteStatusPendingApproval)
	}
	if err := pendingQuote.Approve(); err != nil {
		t.Fatal(err)
	}
	if pendingQuote.Status() != QuoteStatusApproved {
		t.Fatalf("status = %s, want %s", pendingQuote.Status(), QuoteStatusApproved)
	}
}

func TestQuoteRejectsIllegalApprovalTransitions(t *testing.T) {
	quote, err := NewQuote("quote-001", "customer-001")
	if err != nil {
		t.Fatal(err)
	}
	if err := quote.Approve(); !errors.Is(err, ErrQuoteNotApprovable) {
		t.Fatalf("approving draft returned %v", err)
	}
	if err := quote.Reject(); !errors.Is(err, ErrQuoteNotRejectable) {
		t.Fatalf("rejecting draft returned %v", err)
	}
	if err := quote.SubmitForApproval(ApprovalDecision{}); !errors.Is(err, ErrQuoteHasNoLines) {
		t.Fatalf("empty submit returned %v", err)
	}
}
