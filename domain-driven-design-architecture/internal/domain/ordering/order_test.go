package ordering

import (
	"errors"
	"testing"

	"domain-driven-design-architecture/internal/domain/quoting"
)

func approvedQuote(t *testing.T) quoting.Quote {
	t.Helper()
	quote, err := quoting.NewQuote("quote-001", "customer-001")
	if err != nil {
		t.Fatal(err)
	}
	price, err := quoting.NewMoney(14250, "USD")
	if err != nil {
		t.Fatal(err)
	}
	line, err := quoting.NewQuoteLine("sku-001", 2, price)
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

func TestOrderIsCreatedFromApprovedQuoteSnapshot(t *testing.T) {
	order, err := NewOrderFromQuote("order-001", approvedQuote(t))
	if err != nil {
		t.Fatal(err)
	}
	if order.Status() != OrderStatusPendingPayment || order.QuoteID() != "quote-001" || len(order.Lines()) != 1 {
		t.Fatalf("unexpected order %+v", order)
	}
	total, err := order.Total()
	if err != nil {
		t.Fatal(err)
	}
	if total.Cents() != 28500 {
		t.Fatalf("total = %d, want 28500", total.Cents())
	}
}

func TestOrderRequiresApprovedQuote(t *testing.T) {
	quote, err := quoting.NewQuote("quote-001", "customer-001")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := NewOrderFromQuote("order-001", quote); !errors.Is(err, ErrQuoteNotApproved) {
		t.Fatalf("draft conversion returned %v", err)
	}
}

func TestOrderCanBecomePaidOnlyFromPendingPayment(t *testing.T) {
	order, err := NewOrderFromQuote("order-001", approvedQuote(t))
	if err != nil {
		t.Fatal(err)
	}
	if err := order.MarkPaid(); err != nil {
		t.Fatal(err)
	}
	if order.Status() != OrderStatusPaid {
		t.Fatalf("status = %s, want %s", order.Status(), OrderStatusPaid)
	}
	if err := order.MarkPaid(); !errors.Is(err, ErrOrderNotAwaitingPayment) {
		t.Fatalf("repeated paid transition returned %v", err)
	}
}

func TestOrderCanBecomeShippedOnlyAfterPayment(t *testing.T) {
	order, err := NewOrderFromQuote("order-001", approvedQuote(t))
	if err != nil {
		t.Fatal(err)
	}
	if err := order.MarkShipped(); !errors.Is(err, ErrOrderNotShippable) {
		t.Fatalf("unpaid shipment returned %v", err)
	}
	if err := order.MarkPaid(); err != nil {
		t.Fatal(err)
	}
	if err := order.MarkShipped(); err != nil {
		t.Fatal(err)
	}
	if order.Status() != OrderStatusShipped {
		t.Fatalf("status = %s, want %s", order.Status(), OrderStatusShipped)
	}
}

func TestOrderCanBeCancelledBeforeShipment(t *testing.T) {
	order, err := NewOrderFromQuote("order-001", approvedQuote(t))
	if err != nil {
		t.Fatal(err)
	}
	if err := order.MarkPaid(); err != nil {
		t.Fatal(err)
	}
	if err := order.Cancel(); err != nil {
		t.Fatal(err)
	}
	if order.Status() != OrderStatusCancelled {
		t.Fatalf("status = %s, want %s", order.Status(), OrderStatusCancelled)
	}
}

func TestOrderRejectsCancellationAfterShipment(t *testing.T) {
	order, err := NewOrderFromQuote("order-001", approvedQuote(t))
	if err != nil {
		t.Fatal(err)
	}
	if err := order.MarkPaid(); err != nil {
		t.Fatal(err)
	}
	if err := order.MarkShipped(); err != nil {
		t.Fatal(err)
	}
	if err := order.Cancel(); !errors.Is(err, ErrOrderNotCancellable) {
		t.Fatalf("shipped cancellation returned %v", err)
	}
}
