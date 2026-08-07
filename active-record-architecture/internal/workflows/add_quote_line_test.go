package workflows

import (
	"testing"

	"active-record-architecture/internal/records"
)

func setupDraftQuote(t *testing.T, activeProduct bool) (*records.Database, *records.Quote) {
	t.Helper()
	db := records.NewDatabase()
	customer := records.NewCustomer(db, "customer-001", true)
	if err := customer.Save(); err != nil {
		t.Fatalf("Customer.Save() error = %v", err)
	}
	quote, err := records.NewDraftQuote(db, customer.ID)
	if err != nil {
		t.Fatalf("NewDraftQuote() error = %v", err)
	}
	if err := quote.Save(); err != nil {
		t.Fatalf("Quote.Save() error = %v", err)
	}
	product := records.NewProduct(db, "sku-001", "Desk", "Standard", activeProduct, 15000)
	if err := product.Save(); err != nil {
		t.Fatalf("Product.Save() error = %v", err)
	}
	return db, quote
}

func TestAddQuoteLineAppendsAndPersistsProductSnapshot(t *testing.T) {
	db, quote := setupDraftQuote(t, true)

	got, err := AddQuoteLine(db, quote.ID, "sku-001", 2)
	if err != nil {
		t.Fatalf("AddQuoteLine() error = %v", err)
	}
	if len(got.Lines) != 1 {
		t.Fatalf("line count = %d, want 1", len(got.Lines))
	}
	line := got.Lines[0]
	if line.SKU != "sku-001" || line.ProductNameSnapshot != "Desk" || line.ProductCategory != "Standard" {
		t.Fatalf("line snapshot = %#v", line)
	}
	if line.Quantity != 2 || line.LineTotal != 30000 {
		t.Fatalf("line quantities/totals = %#v", line)
	}

	reloaded, err := records.FindQuote(db, quote.ID)
	if err != nil {
		t.Fatalf("FindQuote() error = %v", err)
	}
	if len(reloaded.Lines) != 1 || reloaded.Lines[0] != line {
		t.Fatalf("reloaded lines = %#v, want %#v", reloaded.Lines, got.Lines)
	}
}

func TestAddQuoteLineRejectsInvalidQuantity(t *testing.T) {
	db, quote := setupDraftQuote(t, true)

	_, err := AddQuoteLine(db, quote.ID, "sku-001", 0)
	if err != records.ErrQuantityInvalid {
		t.Fatalf("error = %v, want %v", err, records.ErrQuantityInvalid)
	}
}

func TestAddQuoteLineRejectsUnknownOrInactiveProduct(t *testing.T) {
	db, quote := setupDraftQuote(t, false)

	_, err := AddQuoteLine(db, quote.ID, "sku-404", 1)
	if err != records.ErrProductNotFound {
		t.Fatalf("unknown product error = %v, want %v", err, records.ErrProductNotFound)
	}

	_, err = AddQuoteLine(db, quote.ID, "sku-001", 1)
	if err != records.ErrProductInactive {
		t.Fatalf("inactive product error = %v, want %v", err, records.ErrProductInactive)
	}
}

func TestAddQuoteLineRejectsNonDraftQuote(t *testing.T) {
	db, quote := setupDraftQuote(t, true)
	quote.Status = "Submitted"
	if err := quote.Save(); err != nil {
		t.Fatalf("Quote.Save() error = %v", err)
	}

	_, err := AddQuoteLine(db, quote.ID, "sku-001", 1)
	if err != records.ErrQuoteNotEditable {
		t.Fatalf("error = %v, want %v", err, records.ErrQuoteNotEditable)
	}
}
