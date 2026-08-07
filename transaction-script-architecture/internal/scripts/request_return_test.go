package scripts

import (
	"testing"

	"transaction-script-architecture/internal/data"
)

func TestRequestReturnCreatesRequestedReturnAndRefundState(t *testing.T) {
	store := data.NewStore()
	store.Orders["order-001"] = data.Order{
		ID:     "order-001",
		Status: data.OrderStatusShipped,
		Lines: []data.OrderLine{{
			ID:              "order-001-line-001",
			SKU:             "sku-001",
			ProductCategory: "Standard",
			OrderedQuantity: 2,
			ShippedQuantity: 2,
			UnitPrice:       15000,
		}},
	}

	got, err := RequestReturn(store, "order-001", nil, "Damaged on arrival", "customer-1")
	if err != nil {
		t.Fatalf("RequestReturn() error = %v", err)
	}
	if got.ID != "return-001" || got.Status != data.ReturnStatusRequested {
		t.Fatalf("return request = %#v, want requested return-001", got)
	}
	if got.RequestedBy != "customer-1" {
		t.Fatalf("requested by = %q, want %q", got.RequestedBy, "customer-1")
	}
	if len(got.Lines) != 1 || got.Lines[0].Quantity != 2 {
		t.Fatalf("return lines = %#v, want full shipped quantity", got.Lines)
	}
	if got.RefundStatus != data.RefundStatusNotStarted || got.RefundAmount != 30000 {
		t.Fatalf("refund state = %q amount=%d, want not started / 30000", got.RefundStatus, got.RefundAmount)
	}
	if len(store.Refunds) != 0 {
		t.Fatalf("refund count = %d, want 0 before review", len(store.Refunds))
	}
}

func TestRequestReturnRequiresRequester(t *testing.T) {
	store := data.NewStore()
	store.Orders["order-001"] = data.Order{ID: "order-001", Status: data.OrderStatusShipped, Lines: []data.OrderLine{{
		ID: "line-1", SKU: "sku-001", ShippedQuantity: 1, UnitPrice: 100,
	}}}

	_, err := RequestReturn(store, "order-001", nil, "Missing actor", "")
	if err != ErrActorRequired {
		t.Fatalf("error = %v, want %v", err, ErrActorRequired)
	}
}

func TestRequestReturnRejectsUnshippedOrder(t *testing.T) {
	store := data.NewStore()
	store.Orders["order-001"] = data.Order{ID: "order-001", Status: data.OrderStatusReadyForFulfillment}

	_, err := RequestReturn(store, "order-001", nil, "Changed mind", "customer-1")
	if err != ErrOrderNotReturnable {
		t.Fatalf("error = %v, want %v", err, ErrOrderNotReturnable)
	}
}

func TestRequestReturnRejectsQuantityBeyondShippedQuantity(t *testing.T) {
	store := data.NewStore()
	store.Orders["order-001"] = data.Order{
		ID:     "order-001",
		Status: data.OrderStatusShipped,
		Lines: []data.OrderLine{{
			ID:              "order-001-line-001",
			SKU:             "sku-001",
			ShippedQuantity: 1,
			UnitPrice:       15000,
		}},
	}

	_, err := RequestReturn(store, "order-001", []data.ReturnLine{{
		OrderLineID: "order-001-line-001",
		SKU:         "sku-001",
		Quantity:    2,
	}}, "Too many", "customer-1")
	if err != ErrReturnLinesInvalid {
		t.Fatalf("error = %v, want %v", err, ErrReturnLinesInvalid)
	}
	if len(store.Returns) != 0 {
		t.Fatalf("return count = %d, want 0", len(store.Returns))
	}
}
