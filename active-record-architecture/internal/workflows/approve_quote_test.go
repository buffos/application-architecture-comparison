package workflows

import (
	"testing"

	"active-record-architecture/internal/records"
)

func pendingQuote(t *testing.T) (*records.Database, *records.Quote) {
	t.Helper()
	db, quote := quoteWithLine(t, "CustomBuild")
	if _, err := SubmitQuoteForApproval(db, quote.ID); err != nil {
		t.Fatalf("SubmitQuoteForApproval() error = %v", err)
	}
	return db, quote
}

func TestApproveQuotePersistsApprovalMetadata(t *testing.T) {
	db, quote := pendingQuote(t)

	got, err := ApproveQuote(db, quote.ID, "manager-1", "Approved after review")
	if err != nil {
		t.Fatalf("ApproveQuote() error = %v", err)
	}
	if got.Status != records.QuoteStatusApproved || got.ReviewedBy != "manager-1" || got.DecisionComment != "Approved after review" {
		t.Fatalf("approved quote = %#v", got)
	}

	saved, err := records.FindQuote(db, quote.ID)
	if err != nil {
		t.Fatalf("FindQuote() error = %v", err)
	}
	if saved.Status != records.QuoteStatusApproved || saved.ReviewedBy != "manager-1" {
		t.Fatalf("saved quote = %#v", saved)
	}
}

func TestApproveQuoteRejectsMissingReviewerAndInvalidState(t *testing.T) {
	db, quote := pendingQuote(t)
	if _, err := ApproveQuote(db, quote.ID, "", "Approved"); err != records.ErrReviewerRequired {
		t.Fatalf("missing reviewer error = %v, want %v", err, records.ErrReviewerRequired)
	}

	quote.Status = records.QuoteStatusDraft
	if err := quote.Save(); err != nil {
		t.Fatalf("Quote.Save() error = %v", err)
	}
	if _, err := ApproveQuote(db, quote.ID, "manager-1", "Approved"); err != records.ErrQuoteNotApprovable {
		t.Fatalf("invalid state error = %v, want %v", err, records.ErrQuoteNotApprovable)
	}
}
