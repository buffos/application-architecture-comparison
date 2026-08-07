package scripts

import (
	"testing"

	"transaction-script-architecture/internal/data"
)

func TestGetOrdersAwaitingApprovalListsPendingQuotes(t *testing.T) {
	store := data.NewStore()
	store.Quotes["quote-002"] = data.Quote{
		ID:         "quote-002",
		CustomerID: "customer-002",
		Status:     data.QuoteStatusPendingApproval,
		Lines:      []data.QuoteLine{{ProductCategory: "CustomBuild"}},
	}
	store.Quotes["quote-001"] = data.Quote{
		ID:     "quote-001",
		Status: data.QuoteStatusApproved,
	}
	store.Quotes["quote-003"] = data.Quote{
		ID:     "quote-003",
		Status: data.QuoteStatusPendingApproval,
		Lines:  []data.QuoteLine{{ProductCategory: "Standard"}},
	}

	items, err := GetOrdersAwaitingApproval(store)
	if err != nil {
		t.Fatalf("GetOrdersAwaitingApproval() error = %v", err)
	}
	if len(items) != 2 || items[0].QuoteID != "quote-002" || items[1].QuoteID != "quote-003" {
		t.Fatalf("items = %#v, want sorted pending quotes", items)
	}
	if len(items[0].Reasons) != 1 || items[0].Reasons[0] != ApprovalReasonCustomBuild {
		t.Fatalf("reasons = %#v, want custom-build reason", items[0].Reasons)
	}
}

func TestGetOrdersAwaitingApprovalReturnsEmptyWhenNoPendingQuotes(t *testing.T) {
	items, err := GetOrdersAwaitingApproval(data.NewStore())
	if err != nil {
		t.Fatalf("GetOrdersAwaitingApproval() error = %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("items = %#v, want empty", items)
	}
}
