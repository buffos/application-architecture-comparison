package scripts

import (
	"testing"

	"transaction-script-architecture/internal/data"
)

func TestCreatePartialShipmentLeavesOrderPartiallyShipped(t *testing.T) {
	store := data.NewStore()
	store.Orders["order-001"] = data.Order{
		ID:     "order-001",
		Status: data.OrderStatusReadyForFulfillment,
		Lines: []data.OrderLine{
			{ID: "order-001-line-001", SKU: "sku-001", OrderedQuantity: 2, ReservedQuantity: 2},
			{ID: "order-001-line-002", SKU: "sku-002", OrderedQuantity: 1, ReservedQuantity: 1},
		},
	}
	store.Stocks["sku-001"] = data.StockRecord{SKU: "sku-001", OnHand: 2, Reserved: 2}
	store.Stocks["sku-002"] = data.StockRecord{SKU: "sku-002", OnHand: 1, Reserved: 1}

	shipment, err := CreatePartialShipment(store, "order-001", []data.ShipmentLine{{
		OrderLineID: "order-001-line-001",
		SKU:         "sku-001",
		Quantity:    1,
	}}, "warehouse-1")
	if err != nil {
		t.Fatalf("CreatePartialShipment() error = %v", err)
	}
	if len(shipment.Lines) != 1 || shipment.Lines[0].Quantity != 1 {
		t.Fatalf("shipment lines = %#v, want one unit", shipment.Lines)
	}
	if store.Orders["order-001"].Status != data.OrderStatusPartiallyShipped {
		t.Fatalf("order status = %q, want %q", store.Orders["order-001"].Status, data.OrderStatusPartiallyShipped)
	}
	if stock := store.Stocks["sku-001"]; stock.OnHand != 1 || stock.Reserved != 1 {
		t.Fatalf("sku-001 stock = %#v, want on-hand 1 reserved 1", stock)
	}
	if stock := store.Stocks["sku-002"]; stock.OnHand != 1 || stock.Reserved != 1 {
		t.Fatalf("unshipped stock changed = %#v", stock)
	}
}

func TestCreatePartialShipmentCanFinishOrderLater(t *testing.T) {
	store := data.NewStore()
	store.Orders["order-001"] = data.Order{
		ID:     "order-001",
		Status: data.OrderStatusReadyForFulfillment,
		Lines:  []data.OrderLine{{ID: "line-1", SKU: "sku-001", OrderedQuantity: 2, ReservedQuantity: 2}},
	}
	store.Stocks["sku-001"] = data.StockRecord{SKU: "sku-001", OnHand: 2, Reserved: 2}

	if _, err := CreatePartialShipment(store, "order-001", []data.ShipmentLine{{OrderLineID: "line-1", SKU: "sku-001", Quantity: 1}}, "warehouse-1"); err != nil {
		t.Fatalf("first shipment error = %v", err)
	}
	if _, err := CreateShipment(store, "order-001", "warehouse-1"); err != nil {
		t.Fatalf("finishing shipment error = %v", err)
	}
	if store.Orders["order-001"].Status != data.OrderStatusShipped {
		t.Fatalf("order status = %q, want shipped", store.Orders["order-001"].Status)
	}
	if store.Stocks["sku-001"].OnHand != 0 || store.Stocks["sku-001"].Reserved != 0 {
		t.Fatalf("stock after completion = %#v, want zero", store.Stocks["sku-001"])
	}
}

func TestCreatePartialShipmentRejectsExcessQuantityWithoutMutation(t *testing.T) {
	store := data.NewStore()
	store.Orders["order-001"] = data.Order{
		ID:     "order-001",
		Status: data.OrderStatusReadyForFulfillment,
		Lines:  []data.OrderLine{{ID: "line-1", SKU: "sku-001", OrderedQuantity: 1, ReservedQuantity: 1}},
	}
	store.Stocks["sku-001"] = data.StockRecord{SKU: "sku-001", OnHand: 1, Reserved: 1}

	_, err := CreatePartialShipment(store, "order-001", []data.ShipmentLine{{OrderLineID: "line-1", SKU: "sku-001", Quantity: 2}}, "warehouse-1")
	if err != ErrShipmentLinesInvalid {
		t.Fatalf("error = %v, want %v", err, ErrShipmentLinesInvalid)
	}
	if len(store.Shipments) != 0 || store.Stocks["sku-001"].Reserved != 1 {
		t.Fatalf("mutation after invalid shipment: shipments=%d stock=%#v", len(store.Shipments), store.Stocks["sku-001"])
	}
}
