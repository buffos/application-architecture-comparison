package workflows

import (
	"testing"

	"active-record-architecture/internal/records"
)

func reportQuotes(t *testing.T, converted int, total int) *records.Database {
	t.Helper()
	db := records.NewDatabase()
	customer := records.NewCustomer(db, "customer-001", true)
	if err := customer.Save(); err != nil {
		t.Fatalf("Customer.Save() error = %v", err)
	}
	for index := 0; index < total; index++ {
		quote, err := records.NewDraftQuote(db, customer.ID)
		if err != nil {
			t.Fatalf("NewDraftQuote() error = %v", err)
		}
		if index < converted {
			quote.Status = records.QuoteStatusConverted
			quote.ConvertedOrderID = "order-001"
		}
		if err := quote.Save(); err != nil {
			t.Fatalf("Quote.Save() error = %v", err)
		}
	}
	return db
}

func TestGetQuoteConversionReportCalculatesCountsAndRate(t *testing.T) {
	report, err := records.GetQuoteConversionReport(reportQuotes(t, 1, 2))
	if err != nil {
		t.Fatalf("GetQuoteConversionReport() error = %v", err)
	}
	if report.TotalQuotes != 2 || report.ConvertedQuotes != 1 || report.ConversionRate != 0.5 {
		t.Fatalf("report = %#v, want total 2 converted 1 rate .5", report)
	}
}

func TestGetQuoteConversionReportHandlesEmptyAndCompleteSets(t *testing.T) {
	empty, err := records.GetQuoteConversionReport(records.NewDatabase())
	if err != nil {
		t.Fatalf("empty report error = %v", err)
	}
	if empty.TotalQuotes != 0 || empty.ConvertedQuotes != 0 || empty.ConversionRate != 0 {
		t.Fatalf("empty report = %#v", empty)
	}

	complete, err := records.GetQuoteConversionReport(reportQuotes(t, 3, 3))
	if err != nil {
		t.Fatalf("complete report error = %v", err)
	}
	if complete.TotalQuotes != 3 || complete.ConvertedQuotes != 3 || complete.ConversionRate != 1 {
		t.Fatalf("complete report = %#v", complete)
	}
}

func TestGetQuoteConversionReportRejectsMissingDatabase(t *testing.T) {
	if _, err := records.GetQuoteConversionReport(nil); err != records.ErrDatabaseRequired {
		t.Fatalf("error = %v, want %v", err, records.ErrDatabaseRequired)
	}
}
