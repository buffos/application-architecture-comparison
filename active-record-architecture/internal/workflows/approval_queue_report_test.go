package workflows

import (
	"testing"

	"active-record-architecture/internal/records"
)

func TestGetOrdersAwaitingApprovalListsPendingQuotes(t *testing.T) {
	db, pending := quoteWithLine(t, "CustomBuild")
	if _, err := SubmitQuoteForApproval(db, pending.ID); err != nil {
		t.Fatalf("SubmitQuoteForApproval() pending error = %v", err)
	}

	approved, err := records.NewDraftQuote(db, pending.CustomerID)
	if err != nil {
		t.Fatalf("NewDraftQuote() error = %v", err)
	}
	if err := approved.Save(); err != nil {
		t.Fatalf("approved Quote.Save() error = %v", err)
	}
	standardProduct := records.NewProduct(db, "sku-002", "Standard Desk", "Standard", true, 15000)
	if err := standardProduct.Save(); err != nil {
		t.Fatalf("standard Product.Save() error = %v", err)
	}
	if _, err := AddQuoteLine(db, approved.ID, standardProduct.SKU, 1); err != nil {
		t.Fatalf("AddQuoteLine() approved error = %v", err)
	}
	if _, err := SubmitQuoteForApproval(db, approved.ID); err != nil {
		t.Fatalf("SubmitQuoteForApproval() approved error = %v", err)
	}

	items, err := records.GetOrdersAwaitingApproval(db)
	if err != nil {
		t.Fatalf("GetOrdersAwaitingApproval() error = %v", err)
	}
	if len(items) != 1 || items[0].QuoteID != pending.ID || items[0].CustomerID != pending.CustomerID {
		t.Fatalf("approval queue = %#v, want pending quote only", items)
	}
	if len(items[0].Reasons) != 1 || items[0].Reasons[0] != records.ApprovalReasonCustomBuild {
		t.Fatalf("approval reasons = %#v", items[0].Reasons)
	}
	saved, err := records.FindQuote(db, pending.ID)
	if err != nil {
		t.Fatalf("FindQuote() error = %v", err)
	}
	if saved.Status != records.QuoteStatusPendingApproval {
		t.Fatalf("pending quote after report = %q, want unchanged", saved.Status)
	}
}

func TestGetOrdersAwaitingApprovalReturnsEmptyWhenNoPendingQuotes(t *testing.T) {
	items, err := records.GetOrdersAwaitingApproval(records.NewDatabase())
	if err != nil {
		t.Fatalf("empty approval queue error = %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("empty approval queue = %#v", items)
	}
	if _, err := records.GetOrdersAwaitingApproval(nil); err != records.ErrDatabaseRequired {
		t.Fatalf("missing database error = %v, want %v", err, records.ErrDatabaseRequired)
	}
}
