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
	request, err := RequestReturn(db, order.ID, nil, "damaged on arrival", "customer-1")
	if err != nil {
		t.Fatalf("RequestReturn() error = %v", err)
	}
	return db, request
}

func TestAcceptReturnRecordsReviewWithoutSideEffects(t *testing.T) {
	db, request := requestedReturn(t)

	accepted, err := AcceptReturn(db, request.ID, "reviewer-1", "accept-1")
	if err != nil {
		t.Fatalf("AcceptReturn() error = %v", err)
	}
	if accepted.Status != records.ReturnStatusAccepted || accepted.RefundStatus != records.RefundStatusNotStarted {
		t.Fatalf("accepted return = %#v", accepted)
	}

	order, err := records.FindOrder(db, request.OrderID)
	if err != nil {
		t.Fatalf("FindOrder() error = %v", err)
	}
	if order.Lines[0].ReturnedQuantity != 0 {
		t.Fatalf("returned quantity after review = %d, want 0", order.Lines[0].ReturnedQuantity)
	}
	stock, err := records.FindStock(db, "sku-001")
	if err != nil {
		t.Fatalf("FindStock() error = %v", err)
	}
	if stock.OnHand != 4 || stock.Reserved != 0 {
		t.Fatalf("stock after review = %#v, want on-hand 4 reserved 0", stock)
	}
}

func TestAcceptReturnCannotRepeatReview(t *testing.T) {
	db, request := requestedReturn(t)
	if _, err := AcceptReturn(db, request.ID, "reviewer-1", "accept-1"); err != nil {
		t.Fatalf("first AcceptReturn() error = %v", err)
	}
	if repeated, err := AcceptReturn(db, request.ID, "reviewer-1", "accept-1"); err != nil {
		t.Fatalf("repeated AcceptReturn() error = %v", err)
	} else if repeated.Status != records.ReturnStatusAccepted {
		t.Fatalf("repeated AcceptReturn() result = %#v", repeated)
	}
}
