package shipments

import (
	"errors"
	"testing"
	"time"
)

type fixedClock struct{ now time.Time }

func (c fixedClock) Now() time.Time { return c.now }

func TestShipmentReaderLoadsAndListsShipments(t *testing.T) {
	component := NewComponent(fixedClock{now: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)})
	created, err := component.Create(ShipmentRequest{OrderID: "order-001", CustomerID: "customer-001", Lines: []ShipmentLine{{ProductSKU: "sku-001", ProductName: "Desk", Quantity: 2}}})
	if err != nil {
		t.Fatal(err)
	}
	var reader Reader = component
	details, err := reader.GetShipment(GetShipmentQuery{ShipmentID: created.ID})
	if err != nil {
		t.Fatal(err)
	}
	if details.OrderID != "order-001" || details.LineCount != 1 || details.Lines[0].ProductSKU != "sku-001" {
		t.Fatalf("unexpected details %+v", details)
	}
	listed := reader.ListShipments(ListShipmentsQuery{OrderID: "order-001"})
	if len(listed) != 1 || listed[0].ShipmentID != created.ID {
		t.Fatalf("unexpected list %+v", listed)
	}
}

func TestShipmentReaderRejectsUnknownShipment(t *testing.T) {
	component := NewComponent(fixedClock{})
	_, err := component.GetShipment(GetShipmentQuery{ShipmentID: "shipment-999"})
	if !errors.Is(err, ErrShipmentNotFound) {
		t.Fatalf("got %v", err)
	}
}
