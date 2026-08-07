package workflows

import (
	"testing"

	"active-record-architecture/internal/records"
)

func readyTwoUnitOrder(t *testing.T) (*records.Database, *records.Order) {
	t.Helper()
	db := records.NewDatabase()
	customer := records.NewCustomer(db, "customer-001", true)
	if err := customer.Save(); err != nil {
		t.Fatalf("Customer.Save() error = %v", err)
	}
	product := records.NewProduct(db, "sku-001", "Desk", "Standard", true, 15000)
	if err := product.Save(); err != nil {
		t.Fatalf("Product.Save() error = %v", err)
	}
	quote, err := records.NewDraftQuote(db, customer.ID)
	if err != nil {
		t.Fatalf("NewDraftQuote() error = %v", err)
	}
	if err := quote.Save(); err != nil {
		t.Fatalf("Quote.Save() error = %v", err)
	}
	if _, err := AddQuoteLine(db, quote.ID, product.SKU, 2); err != nil {
		t.Fatalf("AddQuoteLine() error = %v", err)
	}
	if _, err := SubmitQuoteForApproval(db, quote.ID); err != nil {
		t.Fatalf("SubmitQuoteForApproval() error = %v", err)
	}
	stock := records.NewStockRecord(db, product.SKU, 5, 0, 2)
	if err := stock.Save(); err != nil {
		t.Fatalf("StockRecord.Save() error = %v", err)
	}
	order, err := ConvertQuoteToOrder(db, quote.ID, "sales-1")
	if err != nil {
		t.Fatalf("ConvertQuoteToOrder() error = %v", err)
	}
	if _, err := CapturePayment(db, order.ID, records.PaymentOutcomeAccept); err != nil {
		t.Fatalf("CapturePayment() error = %v", err)
	}
	return db, order
}

func TestCreatePartialShipmentRetainsRemainingFulfillment(t *testing.T) {
	db, order := readyTwoUnitOrder(t)
	lines := []records.ShipmentLine{{OrderLineID: order.Lines[0].ID, Quantity: 1}}

	partial, err := CreatePartialShipment(db, order.ID, "warehouse-1", lines)
	if err != nil {
		t.Fatalf("CreatePartialShipment() error = %v", err)
	}
	if len(partial.Lines) != 1 || partial.Lines[0].Quantity != 1 {
		t.Fatalf("partial shipment = %#v", partial)
	}
	saved, err := records.FindOrder(db, order.ID)
	if err != nil {
		t.Fatalf("FindOrder() error = %v", err)
	}
	if saved.Status != records.OrderStatusPartiallyShipped || saved.Lines[0].ShippedQuantity != 1 {
		t.Fatalf("partial order = %#v", saved)
	}
	stock, err := records.FindStock(db, "sku-001")
	if err != nil {
		t.Fatalf("FindStock() error = %v", err)
	}
	if stock.OnHand != 4 || stock.Reserved != 1 {
		t.Fatalf("stock after partial shipment = %#v, want on-hand 4 reserved 1", stock)
	}

	complete, err := CreateShipment(db, order.ID, "warehouse-1")
	if err != nil {
		t.Fatalf("follow-up CreateShipment() error = %v", err)
	}
	if len(complete.Lines) != 1 || complete.Lines[0].Quantity != 1 {
		t.Fatalf("completion shipment = %#v", complete)
	}
	saved, err = records.FindOrder(db, order.ID)
	if err != nil {
		t.Fatalf("FindOrder() after completion error = %v", err)
	}
	if saved.Status != records.OrderStatusShipped || saved.Lines[0].ShippedQuantity != 2 {
		t.Fatalf("completed order = %#v", saved)
	}
	stock, err = records.FindStock(db, "sku-001")
	if err != nil {
		t.Fatalf("FindStock() after completion error = %v", err)
	}
	if stock.OnHand != 3 || stock.Reserved != 0 {
		t.Fatalf("stock after completion = %#v, want on-hand 3 reserved 0", stock)
	}
}

func TestCreatePartialShipmentRejectsInvalidSelectionWithoutMutation(t *testing.T) {
	db, order := readyTwoUnitOrder(t)
	lines := []records.ShipmentLine{{OrderLineID: order.Lines[0].ID, Quantity: 3}}
	if _, err := CreatePartialShipment(db, order.ID, "warehouse-1", lines); err != records.ErrShipmentLinesInvalid {
		t.Fatalf("over-shipment error = %v, want %v", err, records.ErrShipmentLinesInvalid)
	}
	saved, err := records.FindOrder(db, order.ID)
	if err != nil {
		t.Fatalf("FindOrder() error = %v", err)
	}
	if saved.Status != records.OrderStatusReadyForFulfillment || saved.Lines[0].ShippedQuantity != 0 {
		t.Fatalf("order after invalid shipment = %#v", saved)
	}
	stock, err := records.FindStock(db, "sku-001")
	if err != nil {
		t.Fatalf("FindStock() error = %v", err)
	}
	if stock.OnHand != 5 || stock.Reserved != 2 {
		t.Fatalf("stock after invalid shipment = %#v, want on-hand 5 reserved 2", stock)
	}
}
