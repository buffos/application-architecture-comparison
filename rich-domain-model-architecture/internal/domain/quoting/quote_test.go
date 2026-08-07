package quoting

import (
	"errors"
	"testing"
)

func TestQuoteAggregateOwnsBehaviorAndLifecycle(t *testing.T) {
	quote, err := NewQuote("quote-001", "customer-001")
	if err != nil {
		t.Fatal(err)
	}

	price, err := NewMoney(15000, "usd")
	if err != nil {
		t.Fatal(err)
	}
	line, err := NewQuoteLine("sku-001", 2, price)
	if err != nil {
		t.Fatal(err)
	}

	if err := quote.AddLine(line); err != nil {
		t.Fatal(err)
	}
	total, err := quote.Total()
	if err != nil {
		t.Fatal(err)
	}
	if total.Cents() != 30000 || total.Currency() != "USD" {
		t.Fatalf("total = %+v, want 30000 USD", total)
	}

	if err := quote.SubmitForApproval(ApprovalDecision{Required: true}); err != nil {
		t.Fatal(err)
	}
	if quote.Status() != QuoteStatusPendingApproval {
		t.Fatalf("status = %s, want %s", quote.Status(), QuoteStatusPendingApproval)
	}
	if err := quote.Approve(); err != nil {
		t.Fatal(err)
	}
	if quote.Status() != QuoteStatusApproved {
		t.Fatalf("status = %s, want %s", quote.Status(), QuoteStatusApproved)
	}
	if err := quote.AddLine(line); !errors.Is(err, ErrQuoteNotEditable) {
		t.Fatalf("adding after approval returned %v", err)
	}
}

func TestQuoteAggregateRejectsInvalidChanges(t *testing.T) {
	quote, err := NewQuote("quote-001", "customer-001")
	if err != nil {
		t.Fatal(err)
	}
	if err := quote.SubmitForApproval(ApprovalDecision{}); !errors.Is(err, ErrQuoteHasNoLines) {
		t.Fatalf("empty submit returned %v", err)
	}

	price, err := NewMoney(1000, "USD")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := NewQuoteLine("sku-001", 0, price); !errors.Is(err, ErrQuantityMustBePositive) {
		t.Fatalf("invalid quantity returned %v", err)
	}

	line, err := NewQuoteLine("sku-001", 1, price)
	if err != nil {
		t.Fatal(err)
	}
	if err := quote.AddLine(line); err != nil {
		t.Fatal(err)
	}

	euro, err := NewMoney(1000, "EUR")
	if err != nil {
		t.Fatal(err)
	}
	mixed, err := NewQuoteLine("sku-002", 1, euro)
	if err != nil {
		t.Fatal(err)
	}
	if err := quote.AddLine(mixed); !errors.Is(err, ErrCurrencyMismatch) {
		t.Fatalf("mixed currency returned %v", err)
	}
}

func TestQuoteLinesAreExposedAsACopy(t *testing.T) {
	quote, err := NewQuote("quote-001", "customer-001")
	if err != nil {
		t.Fatal(err)
	}
	price, err := NewMoney(1000, "USD")
	if err != nil {
		t.Fatal(err)
	}
	line, err := NewQuoteLine("sku-001", 1, price)
	if err != nil {
		t.Fatal(err)
	}
	if err := quote.AddLine(line); err != nil {
		t.Fatal(err)
	}

	lines := quote.Lines()
	lines[0] = QuoteLine{}
	if quote.Lines()[0].ProductSKU() != "sku-001" {
		t.Fatal("mutating the returned lines changed the aggregate")
	}
}

func TestMoneyProtectsItsValueRules(t *testing.T) {
	if _, err := NewMoney(-1, "USD"); !errors.Is(err, ErrMoneyAmountNegative) {
		t.Fatalf("negative amount returned %v", err)
	}
	if _, err := NewMoney(100, " "); !errors.Is(err, ErrCurrencyRequired) {
		t.Fatalf("missing currency returned %v", err)
	}

	dollars, err := NewMoney(100, "USD")
	if err != nil {
		t.Fatal(err)
	}
	euros, err := NewMoney(100, "EUR")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := dollars.Add(euros); !errors.Is(err, ErrCurrencyMismatch) {
		t.Fatalf("mixed-money addition returned %v", err)
	}
}

func TestQuoteRejectsInvalidLifecycleTransitions(t *testing.T) {
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
}
