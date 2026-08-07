package scripts

import (
	"testing"

	"transaction-script-architecture/internal/data"
)

func TestRejectQuotePersistsRejectionMetadata(t *testing.T) {
	store := data.NewStore()
	store.Quotes["quote-001"] = data.Quote{
		ID:     "quote-001",
		Status: data.QuoteStatusPendingApproval,
	}

	got, err := RejectQuote(store, "quote-001", "manager-1", "Margin exceeds policy")
	if err != nil {
		t.Fatalf("RejectQuote() error = %v", err)
	}

	if got.Status != data.QuoteStatusRejected {
		t.Fatalf("status = %q, want %q", got.Status, data.QuoteStatusRejected)
	}
	if got.ReviewedBy != "manager-1" || got.DecisionComment != "Margin exceeds policy" {
		t.Fatalf("review metadata = %#v, want reviewer and comment", got)
	}

	saved := store.Quotes["quote-001"]
	if saved.Status != data.QuoteStatusRejected {
		t.Fatalf("saved status = %q, want %q", saved.Status, data.QuoteStatusRejected)
	}
}

func TestRejectQuoteRequiresReviewer(t *testing.T) {
	store := data.NewStore()
	store.Quotes["quote-001"] = data.Quote{
		ID:     "quote-001",
		Status: data.QuoteStatusPendingApproval,
	}

	_, err := RejectQuote(store, "quote-001", "", "Rejected")
	if err != ErrReviewerRequired {
		t.Fatalf("error = %v, want %v", err, ErrReviewerRequired)
	}
	if store.Quotes["quote-001"].Status != data.QuoteStatusPendingApproval {
		t.Fatalf("status = %q, want pending approval", store.Quotes["quote-001"].Status)
	}
}

func TestRejectQuoteRejectsNonPendingQuote(t *testing.T) {
	for _, status := range []string{data.QuoteStatusDraft, data.QuoteStatusApproved, data.QuoteStatusRejected} {
		t.Run(status, func(t *testing.T) {
			store := data.NewStore()
			store.Quotes["quote-001"] = data.Quote{
				ID:     "quote-001",
				Status: status,
			}

			_, err := RejectQuote(store, "quote-001", "manager-1", "Rejected")
			if err != ErrQuoteNotRejectable {
				t.Fatalf("error = %v, want %v", err, ErrQuoteNotRejectable)
			}
		})
	}
}
