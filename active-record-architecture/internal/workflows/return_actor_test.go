package workflows

import (
	"testing"

	"active-record-architecture/internal/records"
)

func TestReturnActorsPersistAcrossRequestReviewAndRefund(t *testing.T) {
	db, request := requestedReturn(t)
	if request.RequestedBy != "customer-1" {
		t.Fatalf("requester = %q, want customer-1", request.RequestedBy)
	}

	if _, err := AcceptReturn(db, request.ID, "reviewer-1"); err != nil {
		t.Fatalf("AcceptReturn() error = %v", err)
	}
	if _, err := CompleteRefund(db, request.ID, "processor-1"); err != nil {
		t.Fatalf("CompleteRefund() error = %v", err)
	}

	savedRequest, err := records.FindReturnRequest(db, request.ID)
	if err != nil {
		t.Fatalf("FindReturnRequest() error = %v", err)
	}
	if savedRequest.RequestedBy != "customer-1" || savedRequest.ReviewedBy != "reviewer-1" || savedRequest.ProcessedBy != "processor-1" {
		t.Fatalf("return actors = %#v", savedRequest)
	}
	savedRefund, err := records.FindRefund(db, request.RefundID)
	if err != nil {
		t.Fatalf("FindRefund() error = %v", err)
	}
	if savedRefund.ProcessedBy != "processor-1" {
		t.Fatalf("refund processor = %q, want processor-1", savedRefund.ProcessedBy)
	}
}

func TestRejectedReturnPersistsReviewerWithoutCompletingRefund(t *testing.T) {
	db, request := requestedReturn(t)
	rejected, err := RejectReturn(db, request.ID, "reviewer-2", "not covered")
	if err != nil {
		t.Fatalf("RejectReturn() error = %v", err)
	}
	if rejected.ReviewedBy != "reviewer-2" || rejected.Status != records.ReturnStatusRejected {
		t.Fatalf("rejected return = %#v", rejected)
	}
	refund, err := records.FindRefund(db, request.RefundID)
	if err != nil {
		t.Fatalf("FindRefund() error = %v", err)
	}
	if refund.Status != records.RefundStatusNotStarted || refund.ProcessedBy != "" {
		t.Fatalf("refund after rejection = %#v", refund)
	}
}

func TestMissingReturnActorsAreRejectedBeforeMutation(t *testing.T) {
	db, order := readyToShip(t)
	if _, err := CreateShipment(db, order.ID, "warehouse-1"); err != nil {
		t.Fatalf("CreateShipment() error = %v", err)
	}
	if _, err := RequestReturn(db, order.ID, nil, "missing requester", ""); err != records.ErrActorRequired {
		t.Fatalf("missing requester error = %v, want %v", err, records.ErrActorRequired)
	}
	if _, err := records.FindReturnRequest(db, "return-001"); err != records.ErrReturnNotFound {
		t.Fatalf("return after missing requester = %v, want %v", err, records.ErrReturnNotFound)
	}

	db, request := requestedReturn(t)
	if _, err := AcceptReturn(db, request.ID, ""); err != records.ErrActorRequired {
		t.Fatalf("missing reviewer on accept = %v, want %v", err, records.ErrActorRequired)
	}
	if _, err := RejectReturn(db, request.ID, "", "missing reviewer"); err != records.ErrActorRequired {
		t.Fatalf("missing reviewer on reject = %v, want %v", err, records.ErrActorRequired)
	}
	if _, err := AcceptReturn(db, request.ID, "reviewer-1"); err != nil {
		t.Fatalf("AcceptReturn() error = %v", err)
	}
	if _, err := CompleteRefund(db, request.ID, ""); err != records.ErrActorRequired {
		t.Fatalf("missing processor = %v, want %v", err, records.ErrActorRequired)
	}
	stock, err := records.FindStock(db, "sku-001")
	if err != nil {
		t.Fatalf("FindStock() error = %v", err)
	}
	if stock.OnHand != 4 {
		t.Fatalf("stock after missing processor = %d, want 4", stock.OnHand)
	}
}
