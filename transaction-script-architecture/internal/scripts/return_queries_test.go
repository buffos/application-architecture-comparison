package scripts

import (
	"testing"

	"transaction-script-architecture/internal/data"
)

func TestGetReturnRequestReturnsDefensiveSnapshot(t *testing.T) {
	store := data.NewStore()
	store.Returns["return-001"] = data.ReturnRequest{
		ID:     "return-001",
		Status: data.ReturnStatusRequested,
		Lines:  []data.ReturnLine{{SKU: "sku-001", Quantity: 1}},
	}

	got, err := GetReturnRequest(store, "return-001")
	if err != nil {
		t.Fatalf("GetReturnRequest() error = %v", err)
	}
	got.Lines[0].Quantity = 99
	if store.Returns["return-001"].Lines[0].Quantity != 1 {
		t.Fatalf("store line quantity = %d, want 1 after snapshot mutation", store.Returns["return-001"].Lines[0].Quantity)
	}
}

func TestListReturnRequestsFiltersAndSorts(t *testing.T) {
	store := data.NewStore()
	store.Returns["return-002"] = data.ReturnRequest{ID: "return-002", Status: data.ReturnStatusRejected}
	store.Returns["return-001"] = data.ReturnRequest{ID: "return-001", Status: data.ReturnStatusRequested}
	store.Returns["return-003"] = data.ReturnRequest{ID: "return-003", Status: data.ReturnStatusRequested}

	got, err := ListReturnRequests(store, data.ReturnStatusRequested)
	if err != nil {
		t.Fatalf("ListReturnRequests() error = %v", err)
	}
	if len(got) != 2 || got[0].ID != "return-001" || got[1].ID != "return-003" {
		t.Fatalf("filtered returns = %#v, want return-001 then return-003", got)
	}

	all, err := ListReturnRequests(store, "")
	if err != nil {
		t.Fatalf("ListReturnRequests() all error = %v", err)
	}
	if len(all) != 3 || all[0].ID != "return-001" || all[2].ID != "return-003" {
		t.Fatalf("all returns = %#v, want sorted IDs", all)
	}
}

func TestGetReturnRequestRejectsMissingID(t *testing.T) {
	store := data.NewStore()
	if _, err := GetReturnRequest(store, "return-404"); err != ErrReturnNotFound {
		t.Fatalf("error = %v, want %v", err, ErrReturnNotFound)
	}
}
