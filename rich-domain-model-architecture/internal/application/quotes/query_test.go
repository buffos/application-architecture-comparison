package quotes

import (
	"errors"
	"testing"

	"rich-domain-model-architecture/internal/domain/quoting"
)

func TestQuoteQueryProjectsDetailsAndFiltersStatus(t *testing.T) {
	quote := queryQuote(t)
	reader := NewInMemoryReader()
	if err := reader.Save(quote); err != nil {
		t.Fatal(err)
	}
	details, err := reader.GetQuote(string(quote.ID()))
	if err != nil {
		t.Fatal(err)
	}
	if details.Status != string(quoting.QuoteStatusApproved) || details.TotalCents != 2000 || details.LineCount != 1 {
		t.Fatalf("details = %+v", details)
	}
	rows := reader.ListQuotes(quoting.QuoteStatusApproved)
	if len(rows) != 1 || rows[0].QuoteID != string(quote.ID()) {
		t.Fatalf("rows = %+v", rows)
	}
}

func TestQuoteQueryReturnsNotFound(t *testing.T) {
	if _, err := NewInMemoryReader().GetQuote("missing"); !errors.Is(err, ErrQuoteNotFound) {
		t.Fatalf("missing query returned %v", err)
	}
}

func queryQuote(t *testing.T) quoting.Quote {
	t.Helper()
	price, err := quoting.NewMoney(1000, "USD")
	if err != nil {
		t.Fatal(err)
	}
	line, err := quoting.NewQuoteLine("sku-001", 2, price)
	if err != nil {
		t.Fatal(err)
	}
	quote, err := quoting.NewQuote("quote-query", "customer-001")
	if err != nil {
		t.Fatal(err)
	}
	if err := quote.AddLine(line); err != nil {
		t.Fatal(err)
	}
	if err := quote.SubmitForApproval(quoting.ApprovalDecision{}); err != nil {
		t.Fatal(err)
	}
	return quote
}
