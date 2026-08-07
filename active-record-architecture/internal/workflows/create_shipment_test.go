package workflows

import (
	"testing"

	"active-record-architecture/internal/records"
)

func readyToShip(t *testing.T) (*records.Database, *records.Order) {
	t.Helper()
	db, order := readyOrder(t)
	if _, err := CapturePayment(db, order.ID, records.PaymentOutcomeAccept); err != nil {
		t.Fatalf("CapturePayment() error = %v", err)
	}
	return db, order
}

func TestCreateShipmentConsumesReservedStockAndShipsOrder(t *testing.T) {
	db, order := readyToShip(t)

	shipment, err := CreateShipment(db, order.ID, "warehouse-1")
	if err != nil {
		t.Fatalf("CreateShipment() error = %v", err)
	}
	if shipment.ID != "shipment-001" || shipment.Status != records.ShipmentStatusShipped || len(shipment.Lines) != 1 || shipment.Lines[0].Quantity != 1 {
		t.Fatalf("shipment = %#v", shipment)
	}

	savedOrder, err := records.FindOrder(db, order.ID)
	if err != nil {
		t.Fatalf("FindOrder() error = %v", err)
	}
	if savedOrder.Status != records.OrderStatusShipped || savedOrder.Lines[0].ShippedQuantity != 1 {
		t.Fatalf("saved order = %#v", savedOrder)
	}
	savedShipment, err := records.FindShipment(db, shipment.ID)
	if err != nil {
		t.Fatalf("FindShipment() error = %v", err)
	}
	if savedShipment.Lines[0] != shipment.Lines[0] {
		t.Fatalf("saved shipment line = %#v, want %#v", savedShipment.Lines[0], shipment.Lines[0])
	}
	stock, err := records.FindStock(db, "sku-001")
	if err != nil {
		t.Fatalf("FindStock() error = %v", err)
	}
	if stock.OnHand != 4 || stock.Reserved != 0 {
		t.Fatalf("stock after shipment = %#v, want on-hand 4 reserved 0", stock)
	}
}

func TestCreateShipmentRejectsBeforePaymentAndMissingShipper(t *testing.T) {
	db, order := readyOrder(t)
	if _, err := CreateShipment(db, order.ID, "warehouse-1"); err != records.ErrOrderNotShippable {
		t.Fatalf("payment gate error = %v, want %v", err, records.ErrOrderNotShippable)
	}
	if _, err := CreateShipment(db, order.ID, ""); err != records.ErrShippedByRequired {
		t.Fatalf("shipper error = %v, want %v", err, records.ErrShippedByRequired)
	}
}
