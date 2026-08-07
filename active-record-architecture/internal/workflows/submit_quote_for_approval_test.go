package workflows

import (
	"testing"

	"active-record-architecture/internal/records"
)

func quoteWithLine(t *testing.T, category string) (*records.Database, *records.Quote) {
	t.Helper()
	db := records.NewDatabase()
	customer := records.NewCustomer(db, "customer-001", true)
	if err := customer.Save(); err != nil {
		t.Fatalf("Customer.Save() error = %v", err)
	}
	product := records.NewProduct(db, "sku-001", "Desk", category, true, 15000)
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
	if _, err := AddQuoteLine(db, quote.ID, product.SKU, 1); err != nil {
		t.Fatalf("AddQuoteLine() error = %v", err)
	}
	return db, quote
}

func TestSubmitQuoteForApprovalApprovesStandardQuote(t *testing.T) {
	db, quote := quoteWithLine(t, "Standard")

	got, err := SubmitQuoteForApproval(db, quote.ID)
	if err != nil {
		t.Fatalf("SubmitQuoteForApproval() error = %v", err)
	}
	if got.Status != records.QuoteStatusApproved {
		t.Fatalf("status = %q, want %q", got.Status, records.QuoteStatusApproved)
	}

	saved, err := records.FindQuote(db, quote.ID)
	if err != nil {
		t.Fatalf("FindQuote() error = %v", err)
	}
	if saved.Status != records.QuoteStatusApproved {
		t.Fatalf("saved status = %q, want %q", saved.Status, records.QuoteStatusApproved)
	}
}

func TestSubmitQuoteForApprovalRequiresReviewForCustomBuild(t *testing.T) {
	db, quote := quoteWithLine(t, "CustomBuild")

	got, err := SubmitQuoteForApproval(db, quote.ID)
	if err != nil {
		t.Fatalf("SubmitQuoteForApproval() error = %v", err)
	}
	if got.Status != records.QuoteStatusPendingApproval {
		t.Fatalf("status = %q, want %q", got.Status, records.QuoteStatusPendingApproval)
	}
}

func TestSubmitQuoteForApprovalRejectsEmptyAndNonDraftQuotes(t *testing.T) {
	db := records.NewDatabase()
	customer := records.NewCustomer(db, "customer-001", true)
	if err := customer.Save(); err != nil {
		t.Fatalf("Customer.Save() error = %v", err)
	}
	empty, err := records.NewDraftQuote(db, customer.ID)
	if err != nil {
		t.Fatalf("NewDraftQuote() error = %v", err)
	}
	if err := empty.Save(); err != nil {
		t.Fatalf("Quote.Save() error = %v", err)
	}
	if _, err := SubmitQuoteForApproval(db, empty.ID); err != records.ErrQuoteHasNoLines {
		t.Fatalf("empty quote error = %v, want %v", err, records.ErrQuoteHasNoLines)
	}

	empty.Status = records.QuoteStatusApproved
	if err := empty.Save(); err != nil {
		t.Fatalf("Quote.Save() error = %v", err)
	}
	if _, err := SubmitQuoteForApproval(db, empty.ID); err != records.ErrQuoteNotSubmittable {
		t.Fatalf("non-draft quote error = %v, want %v", err, records.ErrQuoteNotSubmittable)
	}
}
