package workflows

import (
	"testing"
	"time"

	"active-record-architecture/internal/records"
)

func quoteWithReturnWindow(t *testing.T, windowDays int) (*records.Database, *records.Quote) {
	t.Helper()
	db := records.NewDatabase()
	customer := records.NewCustomer(db, "customer-001", true)
	if err := customer.Save(); err != nil {
		t.Fatalf("Customer.Save() error = %v", err)
	}
	product := records.NewProduct(db, "sku-001", "Desk", "Standard", true, 15000)
	product.ReturnWindowDays = windowDays
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
	quote, err = AddQuoteLine(db, quote.ID, product.SKU, 1)
	if err != nil {
		t.Fatalf("AddQuoteLine() error = %v", err)
	}
	return db, quote
}

func TestReturnWindowSnapshotPropagatesToReturnLine(t *testing.T) {
	db, quote := quoteWithReturnWindow(t, 7)
	if quoteLine := quote.Lines[0]; quoteLine.ReturnWindowDays != 7 {
		t.Fatalf("quote window = %d, want 7", quoteLine.ReturnWindowDays)
	}
	if _, err := SubmitQuoteForApproval(db, quote.ID); err != nil {
		t.Fatalf("SubmitQuoteForApproval() error = %v", err)
	}
	stock := records.NewStockRecord(db, "sku-001", 5, 0, 2)
	if err := stock.Save(); err != nil {
		t.Fatalf("StockRecord.Save() error = %v", err)
	}
	order, err := ConvertQuoteToOrder(db, quote.ID, "sales-1")
	if err != nil {
		t.Fatalf("ConvertQuoteToOrder() error = %v", err)
	}
	if order.Lines[0].ReturnWindowDays != 7 {
		t.Fatalf("order window = %d, want 7", order.Lines[0].ReturnWindowDays)
	}
	if _, err := CapturePayment(db, order.ID, records.PaymentOutcomeAccept); err != nil {
		t.Fatalf("CapturePayment() error = %v", err)
	}
	if _, err := CreateShipment(db, order.ID, "warehouse-1"); err != nil {
		t.Fatalf("CreateShipment() error = %v", err)
	}
	request, err := RequestReturnAt(db, order.ID, nil, "window test", "customer-1", time.Date(2026, time.January, 2, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("RequestReturnAt() error = %v", err)
	}
	if request.Lines[0].ReturnWindowDays != 7 || !request.RequestedAt.Equal(time.Date(2026, time.January, 2, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("return snapshot = %#v", request.Lines[0])
	}
}

func TestAcceptReturnAtAllowsInsideWindow(t *testing.T) {
	db, request := requestedReturn(t)
	order, err := records.FindOrder(db, request.OrderID)
	if err != nil {
		t.Fatalf("FindOrder() error = %v", err)
	}
	shippedAt := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	order.ShippedAt = shippedAt
	if err := order.Save(); err != nil {
		t.Fatalf("Order.Save() error = %v", err)
	}

	accepted, err := AcceptReturnAt(db, request.ID, shippedAt.AddDate(0, 0, records.DefaultReturnWindowDays), "reviewer-1", "accept-1")
	if err != nil {
		t.Fatalf("AcceptReturnAt() error = %v", err)
	}
	if accepted.Status != records.ReturnStatusAccepted {
		t.Fatalf("accepted status = %q, want %q", accepted.Status, records.ReturnStatusAccepted)
	}
}

func TestAcceptReturnAtRejectsExpiredWindow(t *testing.T) {
	db, request := requestedReturn(t)
	order, err := records.FindOrder(db, request.OrderID)
	if err != nil {
		t.Fatalf("FindOrder() error = %v", err)
	}
	shippedAt := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	order.ShippedAt = shippedAt
	if err := order.Save(); err != nil {
		t.Fatalf("Order.Save() error = %v", err)
	}

	if _, err := AcceptReturnAt(db, request.ID, shippedAt.AddDate(0, 0, records.DefaultReturnWindowDays+1), "reviewer-1", "accept-1"); err != records.ErrReturnNotEligible {
		t.Fatalf("expired return error = %v, want %v", err, records.ErrReturnNotEligible)
	}
	saved, err := records.FindReturnRequest(db, request.ID)
	if err != nil {
		t.Fatalf("FindReturnRequest() error = %v", err)
	}
	if saved.Status != records.ReturnStatusRequested {
		t.Fatalf("expired request status = %q, want %q", saved.Status, records.ReturnStatusRequested)
	}
}

func TestLegacyReturnLineUsesThirtyDayDefault(t *testing.T) {
	db, request := requestedReturn(t)
	order, err := records.FindOrder(db, request.OrderID)
	if err != nil {
		t.Fatalf("FindOrder() error = %v", err)
	}
	shippedAt := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	order.ShippedAt = shippedAt
	order.Lines[0].ReturnWindowDays = 0
	if err := order.Save(); err != nil {
		t.Fatalf("Order.Save() error = %v", err)
	}
	request.Lines[0].ReturnWindowDays = 0
	if err := request.Save(); err != nil {
		t.Fatalf("ReturnRequest.Save() error = %v", err)
	}

	if _, err := AcceptReturnAt(db, request.ID, shippedAt.AddDate(0, 0, records.DefaultReturnWindowDays+1), "reviewer-1", "accept-1"); err != records.ErrReturnNotEligible {
		t.Fatalf("legacy expired return error = %v, want %v", err, records.ErrReturnNotEligible)
	}
}
