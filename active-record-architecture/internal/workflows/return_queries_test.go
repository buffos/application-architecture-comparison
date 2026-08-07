package workflows

import (
	"testing"

	"active-record-architecture/internal/records"
)

func TestGetReturnRequestReturnsDefensiveSnapshot(t *testing.T) {
	db, request := requestedReturn(t)

	got, err := records.GetReturnRequest(db, request.ID)
	if err != nil {
		t.Fatalf("GetReturnRequest() error = %v", err)
	}
	got.Lines[0].Quantity = 99
	saved, err := records.FindReturnRequest(db, request.ID)
	if err != nil {
		t.Fatalf("FindReturnRequest() error = %v", err)
	}
	if saved.Lines[0].Quantity != 1 {
		t.Fatalf("stored quantity = %d, want 1 after query mutation", saved.Lines[0].Quantity)
	}
}

func TestListReturnRequestsFiltersAndSorts(t *testing.T) {
	db, first := requestedReturn(t)
	second, err := RequestReturn(db, first.OrderID, nil, "second request", "customer-2")
	if err != nil {
		t.Fatalf("second RequestReturn() error = %v", err)
	}
	if _, err := RejectReturn(db, second.ID, "reviewer-1", "not covered", "reject-1"); err != nil {
		t.Fatalf("RejectReturn() error = %v", err)
	}

	requested, err := records.ListReturnRequests(db, records.ReturnStatusRequested)
	if err != nil {
		t.Fatalf("ListReturnRequests() filtered error = %v", err)
	}
	if len(requested) != 1 || requested[0].ID != first.ID {
		t.Fatalf("requested returns = %#v, want only %s", requested, first.ID)
	}

	all, err := records.ListReturnRequests(db, "")
	if err != nil {
		t.Fatalf("ListReturnRequests() all error = %v", err)
	}
	if len(all) != 2 || all[0].ID != "return-001" || all[1].ID != "return-002" {
		t.Fatalf("all returns = %#v, want deterministic ID order", all)
	}
}

func TestGetReturnRequestRejectsMissingID(t *testing.T) {
	db := records.NewDatabase()
	if _, err := records.GetReturnRequest(db, "return-404"); err != records.ErrReturnNotFound {
		t.Fatalf("error = %v, want %v", err, records.ErrReturnNotFound)
	}
}
