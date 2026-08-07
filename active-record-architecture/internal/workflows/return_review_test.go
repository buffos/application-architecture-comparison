package workflows

import (
	"testing"

	"active-record-architecture/internal/records"
)

func TestRejectReturnBlocksRefundAndRestock(t *testing.T) {
	db, request := requestedReturn(t)

	rejected, err := RejectReturn(db, request.ID, "reviewer-1", "outside service policy")
	if err != nil {
		t.Fatalf("RejectReturn() error = %v", err)
	}
	if rejected.Status != records.ReturnStatusRejected || rejected.ReviewNote != "outside service policy" {
		t.Fatalf("rejected return = %#v", rejected)
	}
	if _, err := CompleteRefund(db, request.ID, "processor-1"); err != records.ErrReturnNotRefundable {
		t.Fatalf("rejected refund error = %v, want %v", err, records.ErrReturnNotRefundable)
	}
	stock, err := records.FindStock(db, "sku-001")
	if err != nil {
		t.Fatalf("FindStock() error = %v", err)
	}
	if stock.OnHand != 4 {
		t.Fatalf("on-hand after rejection = %d, want 4", stock.OnHand)
	}
}

func TestCompleteRefundAppliesEffectsOnlyAfterAcceptance(t *testing.T) {
	db, request := requestedReturn(t)
	if _, err := CompleteRefund(db, request.ID, "processor-1"); err != records.ErrReturnNotRefundable {
		t.Fatalf("requested refund error = %v, want %v", err, records.ErrReturnNotRefundable)
	}
	if _, err := AcceptReturn(db, request.ID, "reviewer-1"); err != nil {
		t.Fatalf("AcceptReturn() error = %v", err)
	}

	completed, err := CompleteRefund(db, request.ID, "processor-1")
	if err != nil {
		t.Fatalf("CompleteRefund() error = %v", err)
	}
	if completed.Status != records.ReturnStatusRefunded || completed.RefundStatus != records.RefundStatusCompleted {
		t.Fatalf("completed return = %#v", completed)
	}
	order, err := records.FindOrder(db, request.OrderID)
	if err != nil {
		t.Fatalf("FindOrder() error = %v", err)
	}
	if order.Lines[0].ReturnedQuantity != 1 {
		t.Fatalf("returned quantity = %d, want 1", order.Lines[0].ReturnedQuantity)
	}
	stock, err := records.FindStock(db, "sku-001")
	if err != nil {
		t.Fatalf("FindStock() error = %v", err)
	}
	if stock.OnHand != 5 || stock.Reserved != 0 {
		t.Fatalf("stock after completion = %#v, want on-hand 5 reserved 0", stock)
	}
}

func TestCompleteRefundRejectsInvalidLinesWithoutSideEffects(t *testing.T) {
	db, request := requestedReturn(t)
	if _, err := AcceptReturn(db, request.ID, "reviewer-1"); err != nil {
		t.Fatalf("AcceptReturn() error = %v", err)
	}
	request, err := records.FindReturnRequest(db, request.ID)
	if err != nil {
		t.Fatalf("FindReturnRequest() error = %v", err)
	}
	request.Lines[0].Quantity = 2
	if err := request.Save(); err != nil {
		t.Fatalf("ReturnRequest.Save() error = %v", err)
	}

	if _, err := CompleteRefund(db, request.ID, "processor-1"); err != records.ErrReturnLinesInvalid {
		t.Fatalf("invalid quantity error = %v, want %v", err, records.ErrReturnLinesInvalid)
	}
	order, err := records.FindOrder(db, request.OrderID)
	if err != nil {
		t.Fatalf("FindOrder() error = %v", err)
	}
	if order.Lines[0].ReturnedQuantity != 0 {
		t.Fatalf("returned quantity after rejected completion = %d, want 0", order.Lines[0].ReturnedQuantity)
	}
}

func TestRejectReturnRequiresRequestedState(t *testing.T) {
	db, request := requestedReturn(t)
	if _, err := RejectReturn(db, request.ID, "reviewer-1", "no"); err != nil {
		t.Fatalf("RejectReturn() error = %v", err)
	}
	if _, err := RejectReturn(db, request.ID, "reviewer-1", "again"); err != records.ErrReturnNotRejectable {
		t.Fatalf("repeated RejectReturn() error = %v, want %v", err, records.ErrReturnNotRejectable)
	}
}
