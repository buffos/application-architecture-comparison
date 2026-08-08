package reporting

import (
	"math"
	"testing"

	"rules-engine-architecture/internal/engine"
)

func TestBuildQuoteConversionReportCountsConvertedQuotes(t *testing.T) {
	report := BuildQuoteConversionReport([]engine.QuoteFact{
		{ID: "Q-1001", Status: "Converted"},
		{ID: "Q-1002", Status: "Draft"},
		{ID: "Q-1003", Status: "converted"},
	})

	if report.TotalQuotes != 3 || report.ConvertedQuotes != 2 {
		t.Fatalf("unexpected conversion counts: %+v", report)
	}
	if math.Abs(report.ConversionRatePercent-66.6666667) > 0.0001 {
		t.Fatalf("expected 66.67%% conversion rate, got %f", report.ConversionRatePercent)
	}
}

func TestBuildQuoteConversionReportHandlesEmptyInput(t *testing.T) {
	report := BuildQuoteConversionReport(nil)

	if report.TotalQuotes != 0 || report.ConvertedQuotes != 0 || report.ConversionRatePercent != 0 {
		t.Fatalf("expected empty report, got %+v", report)
	}
}
