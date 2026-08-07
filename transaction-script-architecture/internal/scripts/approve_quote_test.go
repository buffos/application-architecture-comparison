package scripts

import (
	"testing"

	"transaction-script-architecture/internal/data"
)

func TestApproveQuotePersistsApprovalMetadata(t *testing.T) {
	store := data.NewStore()
	store.Quotes["quote-001"] = data.Quote{
		ID:     "quote-001",
		Status: data.QuoteStatusPendingApproval,
	}

	got, err := ApproveQuote(store, "quote-001", "manager-1", "Approved after review")
	if err != nil {
		t.Fatalf("ApproveQuote() error = %v", err)
	}

	if got.Status != data.QuoteStatusApproved {
		t.Fatalf("status = %q, want %q", got.Status, data.QuoteStatusApproved)
	}
	if got.ReviewedBy != "manager-1" {
		t.Fatalf("reviewer = %q, want %q", got.ReviewedBy, "manager-1")
	}
	if got.DecisionComment != "Approved after review" {
		t.Fatalf("decision comment = %q, want %q", got.DecisionComment, "Approved after review")
	}

	saved := store.Quotes["quote-001"]
	if saved.Status != data.QuoteStatusApproved || saved.ReviewedBy != "manager-1" {
		t.Fatalf("saved quote = %#v, want approved quote with reviewer", saved)
	}
}

func TestApproveQuoteRejectsMissingReviewer(t *testing.T) {
	store := data.NewStore()
	store.Quotes["quote-001"] = data.Quote{
		ID:     "quote-001",
		Status: data.QuoteStatusPendingApproval,
	}

	_, err := ApproveQuote(store, "quote-001", "", "Approved")
	if err != ErrReviewerRequired {
		t.Fatalf("error = %v, want %v", err, ErrReviewerRequired)
	}
	if store.Quotes["quote-001"].Status != data.QuoteStatusPendingApproval {
		t.Fatalf("status = %q, want %q", store.Quotes["quote-001"].Status, data.QuoteStatusPendingApproval)
	}
}

func TestApproveQuoteRejectsUnknownQuote(t *testing.T) {
	store := data.NewStore()

	_, err := ApproveQuote(store, "quote-404", "manager-1", "Approved")
	if err != ErrQuoteNotFound {
		t.Fatalf("error = %v, want %v", err, ErrQuoteNotFound)
	}
}

func TestApproveQuoteRejectsNonPendingQuote(t *testing.T) {
	for _, status := range []string{data.QuoteStatusDraft, data.QuoteStatusApproved} {
		t.Run(status, func(t *testing.T) {
			store := data.NewStore()
			store.Quotes["quote-001"] = data.Quote{
				ID:     "quote-001",
				Status: status,
			}

			_, err := ApproveQuote(store, "quote-001", "manager-1", "Approved")
			if err != ErrQuoteNotApprovable {
				t.Fatalf("error = %v, want %v", err, ErrQuoteNotApprovable)
			}
		})
	}
}
