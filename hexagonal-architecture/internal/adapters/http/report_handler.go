package http

import (
	"encoding/json"
	"net/http"

	"hexagonal-architecture/internal/core/application"
)

type ReportHandler struct {
	quoteConversion application.GetQuoteConversionReportUseCase
}

type quoteConversionResponse struct {
	TotalQuotes     int     `json:"totalQuotes"`
	ApprovedQuotes  int     `json:"approvedQuotes"`
	ConvertedQuotes int     `json:"convertedQuotes"`
	ConversionRate  float64 `json:"conversionRate"`
}

func NewReportHandler(quoteConversion application.GetQuoteConversionReportUseCase) ReportHandler {
	return ReportHandler{quoteConversion: quoteConversion}
}

func (h ReportHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet && r.URL.Path == "/reports/quote-conversion" {
		h.quoteConversionReport(w, r)
		return
	}

	http.NotFound(w, r)
}

func (h ReportHandler) quoteConversionReport(w http.ResponseWriter, r *http.Request) {
	report, err := h.quoteConversion.Execute()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(quoteConversionResponse{
		TotalQuotes:     report.TotalQuotes,
		ApprovedQuotes:  report.ApprovedQuotes,
		ConvertedQuotes: report.ConvertedQuotes,
		ConversionRate:  report.ConversionRate,
	})
}
