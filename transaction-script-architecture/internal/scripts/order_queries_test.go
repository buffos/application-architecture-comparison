package scripts

import (
	"testing"

	"transaction-script-architecture/internal/data"
)

func TestGetOrderReturnsDefensiveSnapshot(t *testing.T) {
	store := data.NewStore()
	store.Orders["order-001"] = data.Order{
		ID:     "order-001",
		Status: data.OrderStatusReadyForPayment,
		Lines:  []data.OrderLine{{SKU: "sku-001", OrderedQuantity: 1}},
	}

	got, err := GetOrder(store, "order-001")
	if err != nil {
		t.Fatalf("GetOrder() error = %v", err)
	}
	got.Lines[0].OrderedQuantity = 99
	if store.Orders["order-001"].Lines[0].OrderedQuantity != 1 {
		t.Fatalf("store quantity = %d, want 1 after snapshot mutation", store.Orders["order-001"].Lines[0].OrderedQuantity)
	}
}

func TestListOrdersFiltersAndSorts(t *testing.T) {
	store := data.NewStore()
	store.Orders["order-002"] = data.Order{ID: "order-002", Status: data.OrderStatusShipped}
	store.Orders["order-001"] = data.Order{ID: "order-001", Status: data.OrderStatusReadyForPayment}
	store.Orders["order-003"] = data.Order{ID: "order-003", Status: data.OrderStatusReadyForPayment}

	got, err := ListOrders(store, data.OrderStatusReadyForPayment)
	if err != nil {
		t.Fatalf("ListOrders() error = %v", err)
	}
	if len(got) != 2 || got[0].ID != "order-001" || got[1].ID != "order-003" {
		t.Fatalf("filtered orders = %#v, want order-001 then order-003", got)
	}
}

func TestGetOrderRejectsMissingOrder(t *testing.T) {
	store := data.NewStore()
	if _, err := GetOrder(store, "order-404"); err != ErrOrderNotFound {
		t.Fatalf("error = %v, want %v", err, ErrOrderNotFound)
	}
}
