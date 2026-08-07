package scripts

import (
	"testing"

	"transaction-script-architecture/internal/data"
)

func TestCreateShipmentConsumesReservedStockAndShipsOrder(t *testing.T) {
	store := data.NewStore()
	store.Orders["order-001"] = data.Order{
		ID:     "order-001",
		Status: data.OrderStatusReadyForFulfillment,
		Lines: []data.OrderLine{{
			ID:               "order-001-line-001",
			SKU:              "sku-001",
			OrderedQuantity:  2,
			ReservedQuantity: 2,
		}},
	}
	store.Stocks["sku-001"] = data.StockRecord{SKU: "sku-001", OnHand: 5, Reserved: 2}

	got, err := CreateShipment(store, "order-001", "warehouse-1")
	if err != nil {
		t.Fatalf("CreateShipment() error = %v", err)
	}
	if got.ID != "shipment-001" || got.Status != data.ShipmentStatusShipped {
		t.Fatalf("shipment = %#v, want shipped shipment-001", got)
	}
	if len(got.Lines) != 1 || got.Lines[0].Quantity != 2 {
		t.Fatalf("shipment lines = %#v, want quantity 2", got.Lines)
	}
	if store.Orders["order-001"].Status != data.OrderStatusShipped {
		t.Fatalf("order status = %q, want %q", store.Orders["order-001"].Status, data.OrderStatusShipped)
	}
	if stock := store.Stocks["sku-001"]; stock.OnHand != 3 || stock.Reserved != 0 {
		t.Fatalf("stock after shipment = %#v, want on-hand 3 reserved 0", stock)
	}
}

func TestCreateShipmentRejectsOrderBeforePayment(t *testing.T) {
	store := data.NewStore()
	store.Orders["order-001"] = data.Order{ID: "order-001", Status: data.OrderStatusReadyForPayment}

	_, err := CreateShipment(store, "order-001", "warehouse-1")
	if err != ErrOrderNotShippable {
		t.Fatalf("error = %v, want %v", err, ErrOrderNotShippable)
	}
	if len(store.Shipments) != 0 {
		t.Fatalf("shipment count = %d, want 0", len(store.Shipments))
	}
}

func TestCreateShipmentRequiresShipper(t *testing.T) {
	store := data.NewStore()
	_, err := CreateShipment(store, "order-001", "")
	if err != ErrShippedByRequired {
		t.Fatalf("error = %v, want %v", err, ErrShippedByRequired)
	}
}
