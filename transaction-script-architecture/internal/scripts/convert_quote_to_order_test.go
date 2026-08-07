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
	store.Products["sku-001"] = data.Product{SKU: "sku-001", StockShortagePolicy: data.StockShortageRejectOrder}
	store.Stocks["sku-001"] = data.StockRecord{SKU: "sku-001", OnHand: 5}

	got, err := ConvertQuoteToOrder(store, "quote-001", "sales-1")
	if err != nil {
		t.Fatalf("ConvertQuoteToOrder() error = %v", err)
	}

	if got.ID != "order-001" {
		t.Fatalf("order ID = %q, want %q", got.ID, "order-001")
	}
	if got.Status != data.OrderStatusReadyForPayment {
		t.Fatalf("status = %q, want %q", got.Status, data.OrderStatusReadyForPayment)
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
	if store.Stocks["sku-001"].Reserved != 2 {
		t.Fatalf("reserved stock = %d, want 2", store.Stocks["sku-001"].Reserved)
	}
}

func TestConvertQuoteToOrderRejectsHardStockShortageWithoutMutation(t *testing.T) {
	store := data.NewStore()
	store.Products["sku-001"] = data.Product{SKU: "sku-001", StockShortagePolicy: data.StockShortageRejectOrder}
	store.Stocks["sku-001"] = data.StockRecord{SKU: "sku-001", OnHand: 1}
	store.Quotes["quote-001"] = data.Quote{
		ID:     "quote-001",
		Status: data.QuoteStatusApproved,
		Lines:  []data.QuoteLine{{SKU: "sku-001", Quantity: 2, LineTotal: 10}},
	}

	_, err := ConvertQuoteToOrder(store, "quote-001", "sales-1")
	if err != ErrInsufficientStock {
		t.Fatalf("error = %v, want %v", err, ErrInsufficientStock)
	}
	if len(store.Orders) != 0 {
		t.Fatalf("order count = %d, want 0", len(store.Orders))
	}
	if store.Quotes["quote-001"].Status != data.QuoteStatusApproved {
		t.Fatalf("quote status = %q, want approved", store.Quotes["quote-001"].Status)
	}
	if store.Stocks["sku-001"].Reserved != 0 {
		t.Fatalf("reserved stock = %d, want 0", store.Stocks["sku-001"].Reserved)
	}
}

func TestConvertQuoteToOrderAllowsBackorder(t *testing.T) {
	store := data.NewStore()
	store.Products["sku-001"] = data.Product{SKU: "sku-001", StockShortagePolicy: data.StockShortageAllowBackorder}
	store.Stocks["sku-001"] = data.StockRecord{SKU: "sku-001", OnHand: 1}
	store.Quotes["quote-001"] = data.Quote{
		ID:     "quote-001",
		Status: data.QuoteStatusApproved,
		Lines:  []data.QuoteLine{{SKU: "sku-001", Quantity: 2, LineTotal: 10}},
	}

	got, err := ConvertQuoteToOrder(store, "quote-001", "sales-1")
	if err != nil {
		t.Fatalf("ConvertQuoteToOrder() error = %v", err)
	}
	if got.Status != data.OrderStatusBackordered {
		t.Fatalf("status = %q, want %q", got.Status, data.OrderStatusBackordered)
	}
	if store.Stocks["sku-001"].Reserved != 0 {
		t.Fatalf("reserved stock = %d, want 0 for backorder", store.Stocks["sku-001"].Reserved)
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
