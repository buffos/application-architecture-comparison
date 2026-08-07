package scripts

import (
	"testing"

	"transaction-script-architecture/internal/data"
)

func TestGetShipmentReturnsDefensiveSnapshot(t *testing.T) {
	store := data.NewStore()
	store.Shipments["shipment-001"] = data.Shipment{
		ID:      "shipment-001",
		OrderID: "order-001",
		Lines:   []data.ShipmentLine{{SKU: "sku-001", Quantity: 1}},
	}

	got, err := GetShipment(store, "shipment-001")
	if err != nil {
		t.Fatalf("GetShipment() error = %v", err)
	}
	got.Lines[0].Quantity = 99
	if store.Shipments["shipment-001"].Lines[0].Quantity != 1 {
		t.Fatalf("store quantity = %d, want 1 after snapshot mutation", store.Shipments["shipment-001"].Lines[0].Quantity)
	}
}

func TestListShipmentsFiltersAndSorts(t *testing.T) {
	store := data.NewStore()
	store.Shipments["shipment-002"] = data.Shipment{ID: "shipment-002", OrderID: "order-002"}
	store.Shipments["shipment-001"] = data.Shipment{ID: "shipment-001", OrderID: "order-001"}
	store.Shipments["shipment-003"] = data.Shipment{ID: "shipment-003", OrderID: "order-001"}

	got, err := ListShipments(store, "order-001")
	if err != nil {
		t.Fatalf("ListShipments() error = %v", err)
	}
	if len(got) != 2 || got[0].ID != "shipment-001" || got[1].ID != "shipment-003" {
		t.Fatalf("shipments = %#v, want shipment-001 then shipment-003", got)
	}
}

func TestGetShipmentRejectsMissingShipment(t *testing.T) {
	store := data.NewStore()
	if _, err := GetShipment(store, "shipment-404"); err != ErrShipmentNotFound {
		t.Fatalf("error = %v, want %v", err, ErrShipmentNotFound)
	}
}
