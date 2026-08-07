package workflows

import (
	"testing"

	"active-record-architecture/internal/records"
)

func TestReturnCommandsRequireIdempotencyKeys(t *testing.T) {
	db, request := requestedReturn(t)
	if _, err := AcceptReturn(db, request.ID, "reviewer-1", ""); err != records.ErrIdempotencyKeyRequired {
		t.Fatalf("missing accept key = %v, want %v", err, records.ErrIdempotencyKeyRequired)
	}
	if _, err := RejectReturn(db, request.ID, "reviewer-1", "no", ""); err != records.ErrIdempotencyKeyRequired {
		t.Fatalf("missing reject key = %v, want %v", err, records.ErrIdempotencyKeyRequired)
	}
	if _, err := AcceptReturn(db, request.ID, "reviewer-1", "accept-1"); err != nil {
		t.Fatalf("AcceptReturn() error = %v", err)
	}
	if _, err := CompleteRefund(db, request.ID, "processor-1", ""); err != records.ErrIdempotencyKeyRequired {
		t.Fatalf("missing completion key = %v, want %v", err, records.ErrIdempotencyKeyRequired)
	}
	stock, err := records.FindStock(db, "sku-001")
	if err != nil {
		t.Fatalf("FindStock() error = %v", err)
	}
	if stock.OnHand != 4 {
		t.Fatalf("stock after missing key = %d, want 4", stock.OnHand)
	}
}

func TestCompleteRefundRetryReturnsOriginalAndRestocksOnce(t *testing.T) {
	db, request := requestedReturn(t)
	if _, err := AcceptReturn(db, request.ID, "reviewer-1", "accept-1"); err != nil {
		t.Fatalf("AcceptReturn() error = %v", err)
	}
	first, err := CompleteRefund(db, request.ID, "processor-1", "refund-1")
	if err != nil {
		t.Fatalf("first CompleteRefund() error = %v", err)
	}
	second, err := CompleteRefund(db, request.ID, "processor-2", "refund-1")
	if err != nil {
		t.Fatalf("retry CompleteRefund() error = %v", err)
	}
	if second.ID != first.ID || second.ProcessedBy != "processor-1" || second.Status != records.ReturnStatusRefunded {
		t.Fatalf("retry result = %#v, want original %#v", second, first)
	}
	stock, err := records.FindStock(db, "sku-001")
	if err != nil {
		t.Fatalf("FindStock() error = %v", err)
	}
	if stock.OnHand != 5 || stock.Reserved != 0 {
		t.Fatalf("stock after completion retry = %#v, want on-hand 5 reserved 0", stock)
	}
	if _, err := records.FindRefund(db, "refund-002"); err != records.ErrRefundNotFound {
		t.Fatalf("second refund lookup = %v, want %v", err, records.ErrRefundNotFound)
	}
}
