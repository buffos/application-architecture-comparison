package workflows

import (
	"testing"

	"active-record-architecture/internal/records"
)

func shippedOrder(t *testing.T, db *records.Database, order *records.Order) *records.Shipment {
	t.Helper()
	if _, err := CapturePayment(db, order.ID, records.PaymentOutcomeAccept); err != nil {
		t.Fatalf("CapturePayment() error = %v", err)
	}
	shipment, err := CreateShipment(db, order.ID, "warehouse-1")
	if err != nil {
		t.Fatalf("CreateShipment() error = %v", err)
	}
	return shipment
}

func TestGetShipmentReturnsDefensiveSnapshot(t *testing.T) {
	db, order := readyOrder(t)
	shipment := shippedOrder(t, db, order)

	got, err := records.GetShipment(db, shipment.ID)
	if err != nil {
		t.Fatalf("GetShipment() error = %v", err)
	}
	got.Lines[0].Quantity = 99
	saved, err := records.FindShipment(db, shipment.ID)
	if err != nil {
		t.Fatalf("FindShipment() error = %v", err)
	}
	if saved.Lines[0].Quantity != 1 {
		t.Fatalf("stored shipment quantity = %d, want 1 after query mutation", saved.Lines[0].Quantity)
	}
}

func TestListShipmentsFiltersAndSorts(t *testing.T) {
	db, firstOrder := readyOrder(t)
	first := shippedOrder(t, db, firstOrder)
	secondOrder := secondReadyOrder(t, db)
	second := shippedOrder(t, db, secondOrder)

	filtered, err := records.ListShipments(db, firstOrder.ID)
	if err != nil {
		t.Fatalf("ListShipments() filtered error = %v", err)
	}
	if len(filtered) != 1 || filtered[0].ID != first.ID {
		t.Fatalf("filtered shipments = %#v, want only %s", filtered, first.ID)
	}

	all, err := records.ListShipments(db, "")
	if err != nil {
		t.Fatalf("ListShipments() all error = %v", err)
	}
	if len(all) != 2 || all[0].ID != "shipment-001" || all[1].ID != second.ID {
		t.Fatalf("all shipments = %#v, want deterministic ID order", all)
	}
}

func TestGetShipmentRejectsMissingID(t *testing.T) {
	db := records.NewDatabase()
	if _, err := records.GetShipment(db, "shipment-404"); err != records.ErrShipmentNotFound {
		t.Fatalf("error = %v, want %v", err, records.ErrShipmentNotFound)
	}
}
