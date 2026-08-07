package scripts

import (
	"testing"

	"transaction-script-architecture/internal/data"
)

func TestConvertQuoteToOrderCopiesApprovedQuoteSnapshot(t *testing.T) {
	store := data.NewStore()
	store.Quotes["quote-001"] = data.Quote{
		ID:         "quote-001",
		CustomerID: "customer-001",
		Status:     data.QuoteStatusApproved,
		Lines: []data.QuoteLine{{
			ProductCategory:     "Standard",
			SKU:                 "sku-001",
			ProductNameSnapshot: "Desk",
			Quantity:            2,
			UnitPrice:           15000,
			LineTotal:           30000,
		}},
	}

	got, err := ConvertQuoteToOrder(store, "quote-001", "sales-1")
	if err != nil {
		t.Fatalf("ConvertQuoteToOrder() error = %v", err)
	}

	if got.ID != "order-001" {
		t.Fatalf("order ID = %q, want %q", got.ID, "order-001")
	}
	if got.Status != data.OrderStatusPendingReservation {
		t.Fatalf("status = %q, want %q", got.Status, data.OrderStatusPendingReservation)
	}
	if len(got.Lines) != 1 || got.Lines[0].SKU != "sku-001" || got.Lines[0].LineTotal != 30000 {
		t.Fatalf("order lines = %#v, want copied quote line", got.Lines)
	}
	if store.Quotes["quote-001"].Status != data.QuoteStatusConverted {
		t.Fatalf("quote status = %q, want %q", store.Quotes["quote-001"].Status, data.QuoteStatusConverted)
	}
	if store.Quotes["quote-001"].ConvertedOrderID != got.ID {
		t.Fatalf("converted order ID = %q, want %q", store.Quotes["quote-001"].ConvertedOrderID, got.ID)
	}
}

func TestConvertQuoteToOrderRejectsUnapprovedQuote(t *testing.T) {
	for _, status := range []string{data.QuoteStatusDraft, data.QuoteStatusPendingApproval, data.QuoteStatusRejected} {
		t.Run(status, func(t *testing.T) {
			store := data.NewStore()
			store.Quotes["quote-001"] = data.Quote{ID: "quote-001", Status: status}

			_, err := ConvertQuoteToOrder(store, "quote-001", "sales-1")
			if err != ErrQuoteNotConvertible {
				t.Fatalf("error = %v, want %v", err, ErrQuoteNotConvertible)
			}
			if len(store.Orders) != 0 {
				t.Fatalf("order count = %d, want 0", len(store.Orders))
			}
		})
	}
}

func TestConvertQuoteToOrderRejectsRepeatedConversion(t *testing.T) {
	store := data.NewStore()
	store.Quotes["quote-001"] = data.Quote{ID: "quote-001", Status: data.QuoteStatusConverted}

	_, err := ConvertQuoteToOrder(store, "quote-001", "sales-1")
	if err != ErrQuoteAlreadyConverted {
		t.Fatalf("error = %v, want %v", err, ErrQuoteAlreadyConverted)
	}
}
