package workflows

import (
	"testing"

	"active-record-architecture/internal/records"
)

func requestedReturn(t *testing.T) (*records.Database, *records.ReturnRequest) {
	t.Helper()
	db, order := readyToShip(t)
	if _, err := CreateShipment(db, order.ID, "warehouse-1"); err != nil {
		t.Fatalf("CreateShipment() error = %v", err)
	}
	request, err := RequestReturn(db, order.ID, nil, "damaged on arrival")
	if err != nil {
		t.Fatalf("RequestReturn() error = %v", err)
	}
	return db, request
}

func TestAcceptReturnCompletesRefundAndRestocks(t *testing.T) {
	db, request := requestedReturn(t)

	accepted, err := AcceptReturn(db, request.ID)
	if err != nil {
		t.Fatalf("AcceptReturn() error = %v", err)
	}
	if accepted.Status != records.ReturnStatusRefunded || accepted.RefundStatus != records.RefundStatusCompleted {
		t.Fatalf("accepted return = %#v", accepted)
	}

	order, err := records.FindOrder(db, request.OrderID)
	if err != nil {
		t.Fatalf("FindOrder() error = %v", err)
	}
	if order.Lines[0].ReturnedQuantity != 1 || order.Lines[0].ShippedQuantity != 1 {
		t.Fatalf("order line after return = %#v", order.Lines[0])
	}
	stock, err := records.FindStock(db, "sku-001")
	if err != nil {
		t.Fatalf("FindStock() error = %v", err)
	}
	if stock.OnHand != 5 || stock.Reserved != 0 {
		t.Fatalf("stock after return = %#v, want on-hand 5 reserved 0", stock)
	}
	refund, err := records.FindRefund(db, request.RefundID)
	if err != nil {
		t.Fatalf("FindRefund() error = %v", err)
	}
	if refund.Status != records.RefundStatusCompleted {
		t.Fatalf("refund after return = %#v", refund)
	}
}

func TestAcceptReturnRejectsInvalidLinesWithoutSideEffects(t *testing.T) {
	db, request := requestedReturn(t)
	request.Lines[0].Quantity = 2
	if err := request.Save(); err != nil {
		t.Fatalf("ReturnRequest.Save() error = %v", err)
	}

	if _, err := AcceptReturn(db, request.ID); err != records.ErrReturnLinesInvalid {
		t.Fatalf("invalid quantity error = %v, want %v", err, records.ErrReturnLinesInvalid)
	}
	order, err := records.FindOrder(db, request.OrderID)
	if err != nil {
		t.Fatalf("FindOrder() error = %v", err)
	}
	if order.Lines[0].ReturnedQuantity != 0 {
		t.Fatalf("returned quantity after rejected acceptance = %d, want 0", order.Lines[0].ReturnedQuantity)
	}
	stock, err := records.FindStock(db, "sku-001")
	if err != nil {
		t.Fatalf("FindStock() error = %v", err)
	}
	if stock.OnHand != 4 {
		t.Fatalf("on-hand after rejected acceptance = %d, want 4", stock.OnHand)
	}
}

func TestAcceptReturnCannotRepeatReverseSideEffects(t *testing.T) {
	db, request := requestedReturn(t)
	if _, err := AcceptReturn(db, request.ID); err != nil {
		t.Fatalf("first AcceptReturn() error = %v", err)
	}
	if _, err := AcceptReturn(db, request.ID); err != records.ErrReturnNotAcceptable {
		t.Fatalf("repeated AcceptReturn() error = %v, want %v", err, records.ErrReturnNotAcceptable)
	}
	stock, err := records.FindStock(db, "sku-001")
	if err != nil {
		t.Fatalf("FindStock() error = %v", err)
	}
	if stock.OnHand != 5 {
		t.Fatalf("on-hand after repeated acceptance = %d, want 5", stock.OnHand)
	}
}
