package scripts

import (
	"testing"

	"transaction-script-architecture/internal/data"
)

func TestPartialReturnUpdatesOnlySelectedOrderLine(t *testing.T) {
	store := data.NewStore()
	store.Orders["order-001"] = data.Order{
		ID:     "order-001",
		Status: data.OrderStatusShipped,
		Lines: []data.OrderLine{
			{ID: "line-1", SKU: "sku-001", ProductCategory: "Standard", ShippedQuantity: 2, UnitPrice: 10000},
			{ID: "line-2", SKU: "sku-002", ProductCategory: "Standard", ShippedQuantity: 1, UnitPrice: 20000},
		},
	}
	store.Stocks["sku-001"] = data.StockRecord{SKU: "sku-001", OnHand: 0}
	store.Stocks["sku-002"] = data.StockRecord{SKU: "sku-002", OnHand: 0}

	first, err := RequestReturn(store, "order-001", []data.ReturnLine{{OrderLineID: "line-1", SKU: "sku-001", Quantity: 1}}, "One damaged unit", "customer-1")
	if err != nil {
		t.Fatalf("first RequestReturn() error = %v", err)
	}
	if _, err := AcceptReturn(store, first.ID, "manager-1", "accept-partial-1"); err != nil {
		t.Fatalf("first AcceptReturn() error = %v", err)
	}
	if _, err := CompleteRefund(store, first.ID, "finance-1", "refund-partial-1"); err != nil {
		t.Fatalf("first CompleteRefund() error = %v", err)
	}

	order := store.Orders["order-001"]
	if order.Lines[0].ReturnedQuantity != 1 || order.Lines[1].ReturnedQuantity != 0 {
		t.Fatalf("returned quantities = %#v, want [1 0]", order.Lines)
	}
	if store.Stocks["sku-001"].OnHand != 1 || store.Stocks["sku-002"].OnHand != 0 {
		t.Fatalf("stock after first return = %#v / %#v, want only sku-001 restocked", store.Stocks["sku-001"], store.Stocks["sku-002"])
	}

	second, err := RequestReturn(store, "order-001", []data.ReturnLine{{OrderLineID: "line-1", SKU: "sku-001", Quantity: 1}}, "Remaining damaged unit", "customer-1")
	if err != nil {
		t.Fatalf("second RequestReturn() error = %v", err)
	}
	if _, err := AcceptReturn(store, second.ID, "manager-1", "accept-partial-2"); err != nil {
		t.Fatalf("second AcceptReturn() error = %v", err)
	}
	if _, err := CompleteRefund(store, second.ID, "finance-1", "refund-partial-2"); err != nil {
		t.Fatalf("second CompleteRefund() error = %v", err)
	}
	if store.Orders["order-001"].Lines[0].ReturnedQuantity != 2 {
		t.Fatalf("line-1 returned quantity = %d, want 2", store.Orders["order-001"].Lines[0].ReturnedQuantity)
	}
}

func TestRequestReturnRejectsDuplicateLineEntries(t *testing.T) {
	store := data.NewStore()
	store.Orders["order-001"] = data.Order{
		ID:     "order-001",
		Status: data.OrderStatusShipped,
		Lines:  []data.OrderLine{{ID: "line-1", SKU: "sku-001", ShippedQuantity: 2, UnitPrice: 10000}},
	}

	_, err := RequestReturn(store, "order-001", []data.ReturnLine{
		{OrderLineID: "line-1", SKU: "sku-001", Quantity: 1},
		{OrderLineID: "line-1", SKU: "sku-001", Quantity: 1},
	}, "Duplicate line", "customer-1")
	if err != ErrReturnLinesInvalid {
		t.Fatalf("error = %v, want %v", err, ErrReturnLinesInvalid)
	}
	if len(store.Returns) != 0 {
		t.Fatalf("return count = %d, want 0", len(store.Returns))
	}
}
