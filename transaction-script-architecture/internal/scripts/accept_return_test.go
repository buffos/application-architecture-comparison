package scripts

import (
	"testing"

	"transaction-script-architecture/internal/data"
)

func TestAcceptReturnRecordsReviewBeforeSideEffects(t *testing.T) {
	store := requestedReturnStore()

	got, err := AcceptReturn(store, "return-001", "manager-1", "accept-001")
	if err != nil {
		t.Fatalf("AcceptReturn() error = %v", err)
	}
	if got.Status != data.ReturnStatusAccepted {
		t.Fatalf("status = %q, want %q", got.Status, data.ReturnStatusAccepted)
	}
	if got.ReviewedBy != "manager-1" {
		t.Fatalf("reviewed by = %q, want %q", got.ReviewedBy, "manager-1")
	}
	if store.Stocks["sku-001"].OnHand != 3 {
		t.Fatalf("on-hand stock = %d, want 3 before refund completion", store.Stocks["sku-001"].OnHand)
	}
	if len(store.Refunds) != 0 {
		t.Fatalf("refund count = %d, want 0 before completion", len(store.Refunds))
	}
}

func TestRejectReturnBlocksRefundAndRestock(t *testing.T) {
	store := requestedReturnStore()

	got, err := RejectReturn(store, "return-001", "manager-1", "Product not eligible", "reject-001")
	if err != nil {
		t.Fatalf("RejectReturn() error = %v", err)
	}
	if got.Status != data.ReturnStatusRejected {
		t.Fatalf("status = %q, want %q", got.Status, data.ReturnStatusRejected)
	}
	if got.ReviewedBy != "manager-1" {
		t.Fatalf("reviewed by = %q, want %q", got.ReviewedBy, "manager-1")
	}

	if _, err := CompleteRefund(store, "return-001", "finance-1", "refund-001"); err != ErrReturnNotRefundable {
		t.Fatalf("CompleteRefund() error = %v, want %v", err, ErrReturnNotRefundable)
	}
	if store.Stocks["sku-001"].OnHand != 3 || len(store.Refunds) != 0 {
		t.Fatalf("side effects after rejection: stock=%#v refunds=%d", store.Stocks["sku-001"], len(store.Refunds))
	}
}

func TestCompleteRefundAppliesAcceptedReturn(t *testing.T) {
	store := requestedReturnStore()
	if _, err := AcceptReturn(store, "return-001", "manager-1", "accept-002"); err != nil {
		t.Fatalf("AcceptReturn() error = %v", err)
	}

	got, err := CompleteRefund(store, "return-001", "finance-1", "refund-002")
	if err != nil {
		t.Fatalf("CompleteRefund() error = %v", err)
	}
	if got.Status != data.ReturnStatusRefunded || got.RefundStatus != data.RefundStatusCompleted {
		t.Fatalf("return = %#v, want refunded and completed", got)
	}
	if got.ProcessedBy != "finance-1" || store.Refunds[got.RefundID].ProcessedBy != "finance-1" {
		t.Fatalf("processed by metadata = %#v, want finance-1", got)
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

func TestReturnReviewAndRefundRequireActors(t *testing.T) {
	store := requestedReturnStore()
	if _, err := AcceptReturn(store, "return-001", "", "accept-004"); err != ErrActorRequired {
		t.Fatalf("reviewer error = %v, want %v", err, ErrActorRequired)
	}
	if _, err := RejectReturn(store, "return-001", "", "Rejected", "reject-004"); err != ErrActorRequired {
		t.Fatalf("reject reviewer error = %v, want %v", err, ErrActorRequired)
	}
	request := store.Returns["return-001"]
	request.Status = data.ReturnStatusAccepted
	store.Returns["return-001"] = request
	if _, err := CompleteRefund(store, "return-001", "", "refund-004"); err != ErrActorRequired {
		t.Fatalf("processor error = %v, want %v", err, ErrActorRequired)
	}
}

func TestAcceptReturnRejectsAlreadyProcessedReturn(t *testing.T) {
	store := data.NewStore()
	store.Returns["return-001"] = data.ReturnRequest{ID: "return-001", Status: data.ReturnStatusRefunded}

	_, err := AcceptReturn(store, "return-001", "manager-1", "accept-003")
	if err != ErrReturnNotAcceptable {
		t.Fatalf("error = %v, want %v", err, ErrReturnNotAcceptable)
	}
	if len(store.Refunds) != 0 {
		t.Fatalf("refund count = %d, want 0", len(store.Refunds))
	}
}

func TestCompleteRefundIsIdempotent(t *testing.T) {
	store := requestedReturnStore()
	if _, err := AcceptReturn(store, "return-001", "manager-1", "accept-idempotent"); err != nil {
		t.Fatalf("AcceptReturn() error = %v", err)
	}

	first, err := CompleteRefund(store, "return-001", "finance-1", "refund-idempotent")
	if err != nil {
		t.Fatalf("first CompleteRefund() error = %v", err)
	}
	second, err := CompleteRefund(store, "return-001", "finance-1", "refund-idempotent")
	if err != nil {
		t.Fatalf("duplicate CompleteRefund() error = %v", err)
	}
	if second.ID != first.ID || len(store.Refunds) != 1 {
		t.Fatalf("duplicate result = %#v refunds=%d, want same result and one refund", second, len(store.Refunds))
	}
	if store.Stocks["sku-001"].OnHand != 5 {
		t.Fatalf("on-hand stock = %d, want 5 after one restock", store.Stocks["sku-001"].OnHand)
	}
}

func TestRejectReturnIsIdempotent(t *testing.T) {
	store := requestedReturnStore()
	first, err := RejectReturn(store, "return-001", "manager-1", "Rejected", "reject-idempotent")
	if err != nil {
		t.Fatalf("first RejectReturn() error = %v", err)
	}
	second, err := RejectReturn(store, "return-001", "manager-2", "Different note", "reject-idempotent")
	if err != nil {
		t.Fatalf("duplicate RejectReturn() error = %v", err)
	}
	if second.ID != first.ID || second.ReviewNote != "Rejected" {
		t.Fatalf("duplicate result = %#v, want original rejection", second)
	}
}

func TestReturnCommandsRequireIdempotencyKeys(t *testing.T) {
	store := requestedReturnStore()
	if _, err := AcceptReturn(store, "return-001", "manager-1", ""); err != ErrIdempotencyKeyRequired {
		t.Fatalf("accept key error = %v, want %v", err, ErrIdempotencyKeyRequired)
	}
	if _, err := RejectReturn(store, "return-001", "manager-1", "Rejected", ""); err != ErrIdempotencyKeyRequired {
		t.Fatalf("reject key error = %v, want %v", err, ErrIdempotencyKeyRequired)
	}
	request := store.Returns["return-001"]
	request.Status = data.ReturnStatusAccepted
	store.Returns["return-001"] = request
	if _, err := CompleteRefund(store, "return-001", "finance-1", ""); err != ErrIdempotencyKeyRequired {
		t.Fatalf("refund key error = %v, want %v", err, ErrIdempotencyKeyRequired)
	}
}

func requestedReturnStore() *data.Store {
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
	store.Returns["return-001"] = data.ReturnRequest{
		ID:           "return-001",
		OrderID:      "order-001",
		Status:       data.ReturnStatusRequested,
		RequestedBy:  "customer-1",
		Lines:        []data.ReturnLine{{OrderLineID: "order-001-line-001", SKU: "sku-001", Quantity: 2}},
		RefundStatus: data.RefundStatusNotStarted,
		RefundAmount: 30000,
	}
	return store
}
