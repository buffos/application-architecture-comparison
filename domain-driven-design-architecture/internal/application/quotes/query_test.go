package quotes

import (
	"testing"

	"domain-driven-design-architecture/internal/domain/quoting"
)

func TestReaderProjectsQuoteList(t *testing.T) {
	quote, err := quoting.NewQuote("quote-001", "customer-001")
	if err != nil {
		t.Fatal(err)
	}
	price, err := quoting.NewMoney(1250, "USD")
	if err != nil {
		t.Fatal(err)
	}
	line, err := quoting.NewQuoteLine("sku-001", 2, price)
	if err != nil {
		t.Fatal(err)
	}
	if err := quote.AddLine(line); err != nil {
		t.Fatal(err)
	}
	if err := quote.Submit(); err != nil {
		t.Fatal(err)
	}
	reader := NewInMemoryReader()
	if err := reader.Save(quote); err != nil {
		t.Fatal(err)
	}
	details, err := reader.GetQuote("quote-001")
	if err != nil {
		t.Fatal(err)
	}
	if details.Status != string(quoting.QuoteStatusApproved) || details.TotalCents != 2500 || details.LineCount != 1 {
		t.Fatalf("unexpected details %+v", details)
	}
	if got := reader.ListQuotes(quoting.QuoteStatusApproved); len(got) != 1 || got[0].QuoteID != "quote-001" {
		t.Fatalf("unexpected summaries %+v", got)
	}
}
