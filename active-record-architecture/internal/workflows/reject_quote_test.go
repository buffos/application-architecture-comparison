package workflows

import (
	"testing"

	"active-record-architecture/internal/records"
)

func TestRejectQuotePersistsRejectionMetadata(t *testing.T) {
	db, quote := pendingQuote(t)

	got, err := RejectQuote(db, quote.ID, "manager-2", "Customer credit limit exceeded")
	if err != nil {
		t.Fatalf("RejectQuote() error = %v", err)
	}
	if got.Status != records.QuoteStatusRejected || got.ReviewedBy != "manager-2" || got.DecisionComment != "Customer credit limit exceeded" {
		t.Fatalf("rejected quote = %#v", got)
	}

	saved, err := records.FindQuote(db, quote.ID)
	if err != nil {
		t.Fatalf("FindQuote() error = %v", err)
	}
	if saved.Status != records.QuoteStatusRejected || saved.ReviewedBy != "manager-2" {
		t.Fatalf("saved quote = %#v", saved)
	}
}

func TestRejectQuoteRejectsMissingReviewerAndInvalidState(t *testing.T) {
	db, quote := pendingQuote(t)
	if _, err := RejectQuote(db, quote.ID, "", "Rejected"); err != records.ErrReviewerRequired {
		t.Fatalf("missing reviewer error = %v, want %v", err, records.ErrReviewerRequired)
	}

	quote.Status = records.QuoteStatusApproved
	if err := quote.Save(); err != nil {
		t.Fatalf("Quote.Save() error = %v", err)
	}
	if _, err := RejectQuote(db, quote.ID, "manager-2", "Rejected"); err != records.ErrQuoteNotRejectable {
		t.Fatalf("invalid state error = %v, want %v", err, records.ErrQuoteNotRejectable)
	}
}
