package ordering

import (
	"errors"
	"testing"

	"rich-domain-model-architecture/internal/domain/quoting"
)

func approvedQuote(t *testing.T) quoting.Quote {
	t.Helper()
	price, err := quoting.NewMoney(15000, "USD")
	if err != nil {
		t.Fatal(err)
	}
	line, err := quoting.NewQuoteLineFromProductSnapshotWithCategory("sku-001", "Desk", quoting.ProductCategoryStandard, 2, price)
	if err != nil {
		t.Fatal(err)
	}
	quote, err := quoting.NewQuote("quote-001", "customer-001")
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

func TestOrderIsCreatedOnlyFromApprovedQuote(t *testing.T) {
	quote, err := quoting.NewQuote("quote-001", "customer-001")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := NewOrderFromQuote("order-001", quote); !errors.Is(err, ErrQuoteNotApproved) {
		t.Fatalf("draft conversion returned %v", err)
	}

	quote = approvedQuote(t)
	order, err := NewOrderFromQuote("order-001", quote)
	if err != nil {
		t.Fatal(err)
	}
	if order.Status() != OrderStatusPendingPayment || order.QuoteID() != "quote-001" {
		t.Fatalf("unexpected order state: status=%s quote=%s", order.Status(), order.QuoteID())
	}
	total, err := order.Total()
	if err != nil {
		t.Fatal(err)
	}
	if total.Cents() != 30000 || total.Currency() != "USD" {
		t.Fatalf("total = %+v, want 30000 USD", total)
	}
}

func TestOrderOwnsLineSnapshots(t *testing.T) {
	quote := approvedQuote(t)
	order, err := NewOrderFromQuote("order-001", quote)
	if err != nil {
		t.Fatal(err)
	}
	lines := order.Lines()
	if len(lines) != 1 || lines[0].ProductName() != "Desk" || lines[0].Quantity() != 2 {
		t.Fatalf("unexpected order lines: %+v", lines)
	}
	if err := quote.Reject(); err == nil {
		t.Fatal("approved quote unexpectedly accepted a rejection")
	}
	if order.Status() != OrderStatusPendingPayment {
		t.Fatal("quote operation changed order state")
	}
}
