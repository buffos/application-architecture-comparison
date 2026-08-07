package workflows

import (
	"testing"

	"active-record-architecture/internal/records"
)

func readyTwoLineOrder(t *testing.T) (*records.Database, *records.Order) {
	t.Helper()
	db := records.NewDatabase()
	customer := records.NewCustomer(db, "customer-001", true)
	if err := customer.Save(); err != nil {
		t.Fatalf("Customer.Save() error = %v", err)
	}
	products := []*records.Product{
		records.NewProduct(db, "sku-001", "Desk", "Standard", true, 15000),
		records.NewProduct(db, "sku-002", "Chair", "Standard", true, 7500),
	}
	for _, product := range products {
		if err := product.Save(); err != nil {
			t.Fatalf("Product.Save() error = %v", err)
		}
	}
	quote, err := records.NewDraftQuote(db, customer.ID)
	if err != nil {
		t.Fatalf("NewDraftQuote() error = %v", err)
	}
	if err := quote.Save(); err != nil {
		t.Fatalf("Quote.Save() error = %v", err)
	}
	for _, product := range products {
		if _, err := AddQuoteLine(db, quote.ID, product.SKU, 1); err != nil {
			t.Fatalf("AddQuoteLine() error = %v", err)
		}
		stock := records.NewStockRecord(db, product.SKU, 5, 0, 2)
		if err := stock.Save(); err != nil {
			t.Fatalf("StockRecord.Save() error = %v", err)
		}
	}
	if _, err := SubmitQuoteForApproval(db, quote.ID); err != nil {
		t.Fatalf("SubmitQuoteForApproval() error = %v", err)
	}
	order, err := ConvertQuoteToOrder(db, quote.ID, "sales-1")
	if err != nil {
		t.Fatalf("ConvertQuoteToOrder() error = %v", err)
	}
	if _, err := CapturePayment(db, order.ID, records.PaymentOutcomeAccept); err != nil {
		t.Fatalf("CapturePayment() error = %v", err)
	}
	if _, err := CreateShipment(db, order.ID, "warehouse-1"); err != nil {
		t.Fatalf("CreateShipment() error = %v", err)
	}
	return db, order
}

func TestPartialReturnUpdatesOnlySelectedLineAndThenRemainingQuantity(t *testing.T) {
	db, order := readyTwoLineOrder(t)
	selected := []records.ReturnLine{{OrderLineID: order.Lines[0].ID, Quantity: 1}}
	first, err := RequestReturn(db, order.ID, selected, "desk return", "customer-1")
	if err != nil {
		t.Fatalf("first RequestReturn() error = %v", err)
	}
	if len(first.Lines) != 1 || first.Lines[0].OrderLineID != order.Lines[0].ID {
		t.Fatalf("first return lines = %#v", first.Lines)
	}
	if _, err := AcceptReturn(db, first.ID, "reviewer-1", "accept-1"); err != nil {
		t.Fatalf("first AcceptReturn() error = %v", err)
	}
	if _, err := CompleteRefund(db, first.ID, "processor-1", "refund-1"); err != nil {
		t.Fatalf("first CompleteRefund() error = %v", err)
	}

	saved, err := records.FindOrder(db, order.ID)
	if err != nil {
		t.Fatalf("FindOrder() error = %v", err)
	}
	if saved.Lines[0].ReturnedQuantity != 1 || saved.Lines[1].ReturnedQuantity != 0 {
		t.Fatalf("order after first return = %#v", saved.Lines)
	}
	stock1, err := records.FindStock(db, "sku-001")
	if err != nil {
		t.Fatalf("FindStock(sku-001) error = %v", err)
	}
	stock2, err := records.FindStock(db, "sku-002")
	if err != nil {
		t.Fatalf("FindStock(sku-002) error = %v", err)
	}
	if stock1.OnHand != 5 || stock2.OnHand != 4 {
		t.Fatalf("stock after first return = %#v / %#v", stock1, stock2)
	}

	second, err := RequestReturn(db, order.ID, nil, "chair return", "customer-1")
	if err != nil {
		t.Fatalf("second RequestReturn() error = %v", err)
	}
	if len(second.Lines) != 1 || second.Lines[0].OrderLineID != order.Lines[1].ID || second.Lines[0].Quantity != 1 {
		t.Fatalf("second return lines = %#v, want remaining second line", second.Lines)
	}
	if _, err := AcceptReturn(db, second.ID, "reviewer-1", "accept-2"); err != nil {
		t.Fatalf("second AcceptReturn() error = %v", err)
	}
	if _, err := CompleteRefund(db, second.ID, "processor-1", "refund-2"); err != nil {
		t.Fatalf("second CompleteRefund() error = %v", err)
	}
	saved, err = records.FindOrder(db, order.ID)
	if err != nil {
		t.Fatalf("FindOrder() final error = %v", err)
	}
	if saved.Lines[0].ReturnedQuantity != 1 || saved.Lines[1].ReturnedQuantity != 1 {
		t.Fatalf("final returned quantities = %#v", saved.Lines)
	}
}

func TestRequestReturnRejectsDuplicateLineEntriesBeforePersistence(t *testing.T) {
	db, order := readyTwoLineOrder(t)
	duplicate := []records.ReturnLine{
		{OrderLineID: order.Lines[0].ID, Quantity: 1},
		{OrderLineID: order.Lines[0].ID, Quantity: 1},
	}
	if _, err := RequestReturn(db, order.ID, duplicate, "duplicate line", "customer-1"); err != records.ErrReturnLinesInvalid {
		t.Fatalf("duplicate line error = %v, want %v", err, records.ErrReturnLinesInvalid)
	}
	if _, err := records.FindReturnRequest(db, "return-001"); err != records.ErrReturnNotFound {
		t.Fatalf("return after duplicate rejection = %v, want %v", err, records.ErrReturnNotFound)
	}
}
