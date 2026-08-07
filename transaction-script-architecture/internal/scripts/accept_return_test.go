package scripts

import (
	"testing"

	"transaction-script-architecture/internal/data"
)

func TestAcceptReturnRefundsAndRestocks(t *testing.T) {
	store := data.NewStore()
	store.Orders["order-001"] = data.Order{
		ID:     "order-001",
		Status: data.OrderStatusShipped,
		Lines: []data.OrderLine{{
			ID:              "order-001-line-001",
			SKU:             "sku-001",
			ShippedQuantity: 2,
			UnitPrice:       15000,
		}},
	}
	store.Stocks["sku-001"] = data.StockRecord{SKU: "sku-001", OnHand: 3}

	request, err := RequestReturn(store, "order-001", nil, "Damaged")
	if err != nil {
		t.Fatalf("RequestReturn() error = %v", err)
	}

	got, err := AcceptReturn(store, request.ID)
	if err != nil {
		t.Fatalf("AcceptReturn() error = %v", err)
	}
	if got.Status != data.ReturnStatusRefunded || got.RefundStatus != data.RefundStatusCompleted {
		t.Fatalf("return = %#v, want refunded and completed", got)
	}
	if got.RefundID == "" || store.Refunds[got.RefundID].Amount != 30000 {
		t.Fatalf("refund = %#v, want amount 30000", store.Refunds[got.RefundID])
	}
	if store.Orders["order-001"].Lines[0].ReturnedQuantity != 2 {
		t.Fatalf("returned quantity = %d, want 2", store.Orders["order-001"].Lines[0].ReturnedQuantity)
	}
	if store.Stocks["sku-001"].OnHand != 5 {
		t.Fatalf("on-hand stock = %d, want 5", store.Stocks["sku-001"].OnHand)
	}
}

func TestAcceptReturnRejectsAlreadyProcessedReturn(t *testing.T) {
	store := data.NewStore()
	store.Returns["return-001"] = data.ReturnRequest{ID: "return-001", Status: data.ReturnStatusRefunded}

	_, err := AcceptReturn(store, "return-001")
	if err != ErrReturnNotAcceptable {
		t.Fatalf("error = %v, want %v", err, ErrReturnNotAcceptable)
	}
	if len(store.Refunds) != 0 {
		t.Fatalf("refund count = %d, want 0", len(store.Refunds))
	}
}
