package workflows

import (
	"testing"

	"active-record-architecture/internal/records"
)

func approvedQuote(t *testing.T) (*records.Database, *records.Quote) {
	t.Helper()
	db, quote := quoteWithLine(t, "Standard")
	if _, err := SubmitQuoteForApproval(db, quote.ID); err != nil {
		t.Fatalf("SubmitQuoteForApproval() error = %v", err)
	}
	return db, quote
}

func TestConvertQuoteToOrderCopiesIndependentSnapshot(t *testing.T) {
	db, quote := approvedQuote(t)

	order, err := ConvertQuoteToOrder(db, quote.ID, "sales-1")
	if err != nil {
		t.Fatalf("ConvertQuoteToOrder() error = %v", err)
	}
	if order.ID != "order-001" || order.Status != records.OrderStatusPendingReservation {
		t.Fatalf("order = %#v, want pending order-001", order)
	}
	if len(order.Lines) != 1 || order.Lines[0].SKU != "sku-001" || order.Lines[0].LineTotal != 15000 {
		t.Fatalf("order lines = %#v, want copied quote line", order.Lines)
	}

	savedQuote, err := records.FindQuote(db, quote.ID)
	if err != nil {
		t.Fatalf("FindQuote() error = %v", err)
	}
	if savedQuote.Status != records.QuoteStatusConverted || savedQuote.ConvertedOrderID != order.ID {
		t.Fatalf("saved quote = %#v, want converted quote", savedQuote)
	}
	savedOrder, err := records.FindOrder(db, order.ID)
	if err != nil {
		t.Fatalf("FindOrder() error = %v", err)
	}
	if savedOrder.Lines[0] != order.Lines[0] {
		t.Fatalf("saved order line = %#v, want %#v", savedOrder.Lines[0], order.Lines[0])
	}
}

func TestConvertQuoteToOrderRejectsInvalidSourceState(t *testing.T) {
	db := records.NewDatabase()
	customer := records.NewCustomer(db, "customer-001", true)
	if err := customer.Save(); err != nil {
		t.Fatalf("Customer.Save() error = %v", err)
	}

	for _, status := range []string{records.QuoteStatusDraft, records.QuoteStatusPendingApproval, records.QuoteStatusRejected} {
		t.Run(status, func(t *testing.T) {
			candidate, err := records.NewDraftQuote(db, customer.ID)
			if err != nil {
				t.Fatalf("NewDraftQuote() error = %v", err)
			}
			candidate.Status = status
			if err := candidate.Save(); err != nil {
				t.Fatalf("Quote.Save() error = %v", err)
			}
			if _, err := ConvertQuoteToOrder(db, candidate.ID, "sales-1"); err != records.ErrQuoteNotConvertible {
				t.Fatalf("error = %v, want %v", err, records.ErrQuoteNotConvertible)
			}
		})
	}
}

func TestConvertQuoteToOrderRejectsRepeatedConversion(t *testing.T) {
	db, quote := approvedQuote(t)
	if _, err := ConvertQuoteToOrder(db, quote.ID, "sales-1"); err != nil {
		t.Fatalf("first conversion error = %v", err)
	}
	if _, err := ConvertQuoteToOrder(db, quote.ID, "sales-1"); err != records.ErrQuoteAlreadyConverted {
		t.Fatalf("second conversion error = %v, want %v", err, records.ErrQuoteAlreadyConverted)
	}
}
