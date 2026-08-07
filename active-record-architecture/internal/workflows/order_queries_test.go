package workflows

import (
	"testing"

	"active-record-architecture/internal/records"
)

func secondReadyOrder(t *testing.T, db *records.Database) *records.Order {
	t.Helper()
	quote, err := records.NewDraftQuote(db, "customer-001")
	if err != nil {
		t.Fatalf("NewDraftQuote() error = %v", err)
	}
	if err := quote.Save(); err != nil {
		t.Fatalf("Quote.Save() error = %v", err)
	}
	if _, err := AddQuoteLine(db, quote.ID, "sku-001", 1); err != nil {
		t.Fatalf("AddQuoteLine() error = %v", err)
	}
	if _, err := SubmitQuoteForApproval(db, quote.ID); err != nil {
		t.Fatalf("SubmitQuoteForApproval() error = %v", err)
	}
	order, err := ConvertQuoteToOrder(db, quote.ID, "sales-2")
	if err != nil {
		t.Fatalf("ConvertQuoteToOrder() error = %v", err)
	}
	return order
}

func TestGetOrderReturnsDefensiveSnapshot(t *testing.T) {
	db, order := readyOrder(t)

	got, err := records.GetOrder(db, order.ID)
	if err != nil {
		t.Fatalf("GetOrder() error = %v", err)
	}
	got.Lines[0].OrderedQuantity = 99
	saved, err := records.FindOrder(db, order.ID)
	if err != nil {
		t.Fatalf("FindOrder() error = %v", err)
	}
	if saved.Lines[0].OrderedQuantity != 1 {
		t.Fatalf("stored quantity = %d, want 1 after query mutation", saved.Lines[0].OrderedQuantity)
	}
}

func TestListOrdersFiltersAndSorts(t *testing.T) {
	db, first := readyOrder(t)
	second := secondReadyOrder(t, db)
	if _, err := CapturePayment(db, second.ID, records.PaymentOutcomeAccept); err != nil {
		t.Fatalf("CapturePayment() error = %v", err)
	}
	if _, err := CreateShipment(db, second.ID, "warehouse-1"); err != nil {
		t.Fatalf("CreateShipment() error = %v", err)
	}

	ready, err := records.ListOrders(db, records.OrderStatusReadyForPayment)
	if err != nil {
		t.Fatalf("ListOrders() filtered error = %v", err)
	}
	if len(ready) != 1 || ready[0].ID != first.ID {
		t.Fatalf("ready orders = %#v, want only %s", ready, first.ID)
	}

	all, err := records.ListOrders(db, "")
	if err != nil {
		t.Fatalf("ListOrders() all error = %v", err)
	}
	if len(all) != 2 || all[0].ID != "order-001" || all[1].ID != "order-002" {
		t.Fatalf("all orders = %#v, want deterministic ID order", all)
	}
}

func TestGetOrderRejectsMissingID(t *testing.T) {
	db := records.NewDatabase()
	if _, err := records.GetOrder(db, "order-404"); err != records.ErrOrderNotFound {
		t.Fatalf("error = %v, want %v", err, records.ErrOrderNotFound)
	}
}
