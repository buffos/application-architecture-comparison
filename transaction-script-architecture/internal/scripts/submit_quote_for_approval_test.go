package scripts

import (
	"testing"

	"transaction-script-architecture/internal/data"
)

func TestSubmitQuoteForApprovalApprovesStandardQuote(t *testing.T) {
	store := data.NewStore()
	store.Quotes["quote-001"] = data.Quote{
		ID:     "quote-001",
		Status: data.QuoteStatusDraft,
		Lines: []data.QuoteLine{
			{ProductCategory: "Standard", SKU: "sku-001", Quantity: 1},
		},
	}

	got, err := SubmitQuoteForApproval(store, "quote-001")
	if err != nil {
		t.Fatalf("SubmitQuoteForApproval() error = %v", err)
	}
	if got.Status != data.QuoteStatusApproved {
		t.Fatalf("status = %q, want %q", got.Status, data.QuoteStatusApproved)
	}
	if store.Quotes["quote-001"].Status != data.QuoteStatusApproved {
		t.Fatalf("saved status = %q, want %q", store.Quotes["quote-001"].Status, data.QuoteStatusApproved)
	}
}

func TestSubmitQuoteForApprovalRequiresApprovalForCustomBuild(t *testing.T) {
	store := data.NewStore()
	store.Quotes["quote-001"] = data.Quote{
		ID:     "quote-001",
		Status: data.QuoteStatusDraft,
		Lines: []data.QuoteLine{
			{ProductCategory: "CustomBuild", SKU: "sku-001", Quantity: 1},
		},
	}

	got, err := SubmitQuoteForApproval(store, "quote-001")
	if err != nil {
		t.Fatalf("SubmitQuoteForApproval() error = %v", err)
	}
	if got.Status != data.QuoteStatusPendingApproval {
		t.Fatalf("status = %q, want %q", got.Status, data.QuoteStatusPendingApproval)
	}
}

func TestSubmitQuoteForApprovalRejectsEmptyQuote(t *testing.T) {
	store := data.NewStore()
	store.Quotes["quote-001"] = data.Quote{
		ID:     "quote-001",
		Status: data.QuoteStatusDraft,
	}

	_, err := SubmitQuoteForApproval(store, "quote-001")
	if err != ErrQuoteHasNoLines {
		t.Fatalf("error = %v, want %v", err, ErrQuoteHasNoLines)
	}
}

func TestSubmitQuoteForApprovalRejectsNonDraftQuote(t *testing.T) {
	store := data.NewStore()
	store.Quotes["quote-001"] = data.Quote{
		ID:     "quote-001",
		Status: data.QuoteStatusApproved,
		Lines: []data.QuoteLine{
			{ProductCategory: "Standard", SKU: "sku-001", Quantity: 1},
		},
	}

	_, err := SubmitQuoteForApproval(store, "quote-001")
	if err != ErrQuoteNotSubmittable {
		t.Fatalf("error = %v, want %v", err, ErrQuoteNotSubmittable)
	}
}
