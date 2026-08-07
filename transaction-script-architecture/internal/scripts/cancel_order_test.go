package scripts

import (
	"testing"

	"transaction-script-architecture/internal/data"
)

func TestCancelOrderReleasesReservation(t *testing.T) {
	store := data.NewStore()
	store.Orders["order-001"] = data.Order{
		ID:     "order-001",
		Status: data.OrderStatusReadyForPayment,
		Lines:  []data.OrderLine{{SKU: "sku-001", ReservedQuantity: 2}},
	}
	store.Stocks["sku-001"] = data.StockRecord{SKU: "sku-001", OnHand: 5, Reserved: 2}

	got, err := CancelOrder(store, "order-001", "sales-1", "Customer withdrew request")
	if err != nil {
		t.Fatalf("CancelOrder() error = %v", err)
	}
	if got.Status != data.OrderStatusCancelled {
		t.Fatalf("status = %q, want %q", got.Status, data.OrderStatusCancelled)
	}
	if got.CancelledBy != "sales-1" || got.CancellationReason != "Customer withdrew request" {
		t.Fatalf("cancellation metadata = %#v, want actor and reason", got)
	}
	if stock := store.Stocks["sku-001"]; stock.Reserved != 0 {
		t.Fatalf("reserved stock = %d, want 0", stock.Reserved)
	}
	if got.Lines[0].ReservedQuantity != 0 {
		t.Fatalf("line reserved quantity = %d, want 0", got.Lines[0].ReservedQuantity)
	}
}

func TestCancelOrderRejectsShippedOrder(t *testing.T) {
	store := data.NewStore()
	store.Orders["order-001"] = data.Order{ID: "order-001", Status: data.OrderStatusShipped}

	_, err := CancelOrder(store, "order-001", "sales-1", "Too late")
	if err != ErrOrderNotCancellable {
		t.Fatalf("error = %v, want %v", err, ErrOrderNotCancellable)
	}
}

func TestCancelOrderRequiresActorAndReason(t *testing.T) {
	store := data.NewStore()
	store.Orders["order-001"] = data.Order{ID: "order-001", Status: data.OrderStatusReadyForPayment}

	if _, err := CancelOrder(store, "order-001", "", "reason"); err != ErrCancelledByRequired {
		t.Fatalf("actor error = %v, want %v", err, ErrCancelledByRequired)
	}
	if _, err := CancelOrder(store, "order-001", "sales-1", ""); err != ErrCancellationReasonRequired {
		t.Fatalf("reason error = %v, want %v", err, ErrCancellationReasonRequired)
	}
}
