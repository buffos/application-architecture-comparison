package workflows

import (
	"testing"

	"active-record-architecture/internal/records"
)

func TestGetQuoteReturnsDefensiveSnapshot(t *testing.T) {
	db, quote := quoteWithLine(t, "Standard")

	got, err := records.GetQuote(db, quote.ID)
	if err != nil {
		t.Fatalf("GetQuote() error = %v", err)
	}
	got.Lines[0].Quantity = 99
	saved, err := records.FindQuote(db, quote.ID)
	if err != nil {
		t.Fatalf("FindQuote() error = %v", err)
	}
	if saved.Lines[0].Quantity != 1 {
		t.Fatalf("stored quantity = %d, want 1 after query mutation", saved.Lines[0].Quantity)
	}
}

func TestListQuotesFiltersByStatusAndCustomer(t *testing.T) {
	db, first := quoteWithLine(t, "Standard")
	first.Status = records.QuoteStatusApproved
	if err := first.Save(); err != nil {
		t.Fatalf("Quote.Save() error = %v", err)
	}
	secondCustomer := records.NewCustomer(db, "customer-002", true)
	if err := secondCustomer.Save(); err != nil {
		t.Fatalf("Customer.Save() error = %v", err)
	}
	second, err := records.NewDraftQuote(db, secondCustomer.ID)
	if err != nil {
		t.Fatalf("NewDraftQuote() error = %v", err)
	}
	second.Status = records.QuoteStatusApproved
	if err := second.Save(); err != nil {
		t.Fatalf("second Quote.Save() error = %v", err)
	}
	third, err := records.NewDraftQuote(db, first.CustomerID)
	if err != nil {
		t.Fatalf("third NewDraftQuote() error = %v", err)
	}
	if err := third.Save(); err != nil {
		t.Fatalf("third Quote.Save() error = %v", err)
	}

	filtered, err := records.ListQuotes(db, records.QuoteStatusApproved, first.CustomerID)
	if err != nil {
		t.Fatalf("ListQuotes() filtered error = %v", err)
	}
	if len(filtered) != 1 || filtered[0].ID != first.ID {
		t.Fatalf("filtered quotes = %#v, want only %s", filtered, first.ID)
	}

	byStatus, err := records.ListQuotes(db, records.QuoteStatusApproved, "")
	if err != nil {
		t.Fatalf("ListQuotes() status error = %v", err)
	}
	if len(byStatus) != 2 || byStatus[0].ID != "quote-001" || byStatus[1].ID != "quote-002" {
		t.Fatalf("status quotes = %#v, want sorted approved quotes", byStatus)
	}
}

func TestGetQuoteRejectsMissingID(t *testing.T) {
	db := records.NewDatabase()
	if _, err := records.GetQuote(db, "quote-404"); err != records.ErrQuoteNotFound {
		t.Fatalf("error = %v, want %v", err, records.ErrQuoteNotFound)
	}
}
