package workflows

import (
	"testing"

	"active-record-architecture/internal/records"
)

func TestCancelOrderReleasesReservationAndPersistsMetadata(t *testing.T) {
	db, order := readyOrder(t)

	cancelled, err := CancelOrder(db, order.ID, "sales-1", "customer changed mind")
	if err != nil {
		t.Fatalf("CancelOrder() error = %v", err)
	}
	if cancelled.Status != records.OrderStatusCancelled || cancelled.CancelledBy != "sales-1" || cancelled.CancellationReason != "customer changed mind" {
		t.Fatalf("cancelled order = %#v", cancelled)
	}
	if cancelled.Lines[0].ReservedQuantity != 0 {
		t.Fatalf("cancelled line reservation = %d, want 0", cancelled.Lines[0].ReservedQuantity)
	}

	savedOrder, err := records.FindOrder(db, order.ID)
	if err != nil {
		t.Fatalf("FindOrder() error = %v", err)
	}
	if savedOrder.Status != records.OrderStatusCancelled || savedOrder.CancelledBy != "sales-1" {
		t.Fatalf("saved order = %#v", savedOrder)
	}
	stock, err := records.FindStock(db, "sku-001")
	if err != nil {
		t.Fatalf("FindStock() error = %v", err)
	}
	if stock.Reserved != 0 || stock.OnHand != 5 {
		t.Fatalf("stock after cancellation = %#v, want reserved 0 on-hand 5", stock)
	}
}

func TestCancelOrderRejectsRepeatedOrShippedCancellation(t *testing.T) {
	db, order := readyOrder(t)
	if _, err := CancelOrder(db, order.ID, "sales-1", "duplicate request"); err != nil {
		t.Fatalf("first CancelOrder() error = %v", err)
	}
	if _, err := CancelOrder(db, order.ID, "sales-1", "duplicate request"); err != records.ErrOrderNotCancellable {
		t.Fatalf("repeated cancellation error = %v, want %v", err, records.ErrOrderNotCancellable)
	}

	db, order = readyToShip(t)
	if _, err := CreateShipment(db, order.ID, "warehouse-1"); err != nil {
		t.Fatalf("CreateShipment() error = %v", err)
	}
	if _, err := CancelOrder(db, order.ID, "sales-1", "too late"); err != records.ErrOrderNotCancellable {
		t.Fatalf("shipped cancellation error = %v, want %v", err, records.ErrOrderNotCancellable)
	}
}

func TestCancelOrderValidatesBeforeChangingState(t *testing.T) {
	db, order := readyOrder(t)
	if _, err := CancelOrder(db, order.ID, "", "missing actor"); err != records.ErrCancelledByRequired {
		t.Fatalf("actor error = %v, want %v", err, records.ErrCancelledByRequired)
	}
	if _, err := CancelOrder(db, order.ID, "sales-1", ""); err != records.ErrCancellationReasonRequired {
		t.Fatalf("reason error = %v, want %v", err, records.ErrCancellationReasonRequired)
	}
	stock, err := records.FindStock(db, "sku-001")
	if err != nil {
		t.Fatalf("FindStock() error = %v", err)
	}
	if stock.Reserved != 1 {
		t.Fatalf("stock after rejected cancellation = %#v, want reservation 1", stock)
	}
}

func TestCancelOrderRejectsInvalidReservationWithoutPartialRelease(t *testing.T) {
	db, order := readyOrder(t)
	stock, err := records.FindStock(db, "sku-001")
	if err != nil {
		t.Fatalf("FindStock() error = %v", err)
	}
	stock.Reserved = 0
	if err := stock.Save(); err != nil {
		t.Fatalf("StockRecord.Save() error = %v", err)
	}

	if _, err := CancelOrder(db, order.ID, "sales-1", "invalid reservation"); err != records.ErrStockReleaseInvalid {
		t.Fatalf("reservation error = %v, want %v", err, records.ErrStockReleaseInvalid)
	}
	savedOrder, err := records.FindOrder(db, order.ID)
	if err != nil {
		t.Fatalf("FindOrder() error = %v", err)
	}
	if savedOrder.Status == records.OrderStatusCancelled || savedOrder.Lines[0].ReservedQuantity != 1 {
		t.Fatalf("saved order after failed cancellation = %#v", savedOrder)
	}
}
