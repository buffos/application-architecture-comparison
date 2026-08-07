package scripts

import (
	"testing"

	"transaction-script-architecture/internal/data"
)

func TestGetQuoteConversionReportCalculatesCountsAndRate(t *testing.T) {
	store := data.NewStore()
	store.Quotes["quote-001"] = data.Quote{ID: "quote-001", Status: data.QuoteStatusConverted}
	store.Quotes["quote-002"] = data.Quote{ID: "quote-002", Status: data.QuoteStatusApproved}
	store.Quotes["quote-003"] = data.Quote{ID: "quote-003", Status: data.QuoteStatusDraft}

	report, err := GetQuoteConversionReport(store)
	if err != nil {
		t.Fatalf("GetQuoteConversionReport() error = %v", err)
	}
	if report.TotalQuotes != 3 || report.ConvertedQuotes != 1 {
		t.Fatalf("report = %#v, want total 3 converted 1", report)
	}
	if report.ConversionRate < 0.333 || report.ConversionRate > 0.334 {
		t.Fatalf("conversion rate = %f, want about 0.333", report.ConversionRate)
	}
}

func TestGetQuoteConversionReportHandlesEmptyStore(t *testing.T) {
	report, err := GetQuoteConversionReport(data.NewStore())
	if err != nil {
		t.Fatalf("GetQuoteConversionReport() error = %v", err)
	}
	if report.TotalQuotes != 0 || report.ConvertedQuotes != 0 || report.ConversionRate != 0 {
		t.Fatalf("report = %#v, want all zero", report)
	}
}
