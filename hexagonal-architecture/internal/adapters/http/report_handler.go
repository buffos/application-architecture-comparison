package http

import (
	"encoding/json"
	"net/http"

	"hexagonal-architecture/internal/core/application"
)

type ReportHandler struct {
	quoteConversion       application.GetQuoteConversionReportUseCase
	returnRateByCategory  application.GetReturnRateByCategoryReportUseCase
	topDiscountedProducts application.GetTopDiscountedProductsReportUseCase
	lowStockItems         application.GetLowStockItemsReportUseCase
}

type quoteConversionResponse struct {
	TotalQuotes     int     `json:"totalQuotes"`
	ApprovedQuotes  int     `json:"approvedQuotes"`
	ConvertedQuotes int     `json:"convertedQuotes"`
	ConversionRate  float64 `json:"conversionRate"`
}

type returnRateByCategoryResponse struct {
	Categories []returnRateByCategoryRowResponse `json:"categories"`
}

type returnRateByCategoryRowResponse struct {
	Category        string  `json:"category"`
	ShippedQuantity int     `json:"shippedQuantity"`
	ReturnQuantity  int     `json:"returnQuantity"`
	ReturnRate      float64 `json:"returnRate"`
}

type topDiscountedProductsResponse struct {
	Products []topDiscountedProductRowResponse `json:"products"`
}

type topDiscountedProductRowResponse struct {
	SKU                 string  `json:"sku"`
	ProductName         string  `json:"productName"`
	QuotedQuantity      int     `json:"quotedQuantity"`
	TotalDiscountAmount int     `json:"totalDiscountAmount"`
	AverageDiscountRate float64 `json:"averageDiscountRate"`
}

type lowStockItemsResponse struct {
	Items []lowStockItemRowResponse `json:"items"`
}

type lowStockItemRowResponse struct {
	SKU              string `json:"sku"`
	Available        int    `json:"available"`
	ReorderThreshold int    `json:"reorderThreshold"`
}

func NewReportHandler(
	quoteConversion application.GetQuoteConversionReportUseCase,
	returnRateByCategory application.GetReturnRateByCategoryReportUseCase,
	topDiscountedProducts application.GetTopDiscountedProductsReportUseCase,
	lowStockItems application.GetLowStockItemsReportUseCase,
) ReportHandler {
	return ReportHandler{
		quoteConversion:       quoteConversion,
		returnRateByCategory:  returnRateByCategory,
		topDiscountedProducts: topDiscountedProducts,
		lowStockItems:         lowStockItems,
	}
}

func (h ReportHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet && r.URL.Path == "/reports/quote-conversion" {
		h.quoteConversionReport(w, r)
		return
	}
	if r.Method == http.MethodGet && r.URL.Path == "/reports/return-rate-by-category" {
		h.returnRateByCategoryReport(w, r)
		return
	}
	if r.Method == http.MethodGet && r.URL.Path == "/reports/top-discounted-products" {
		h.topDiscountedProductsReport(w, r)
		return
	}
	if r.Method == http.MethodGet && r.URL.Path == "/reports/low-stock-items" {
		h.lowStockItemsReport(w, r)
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

func (h ReportHandler) returnRateByCategoryReport(w http.ResponseWriter, r *http.Request) {
	report, err := h.returnRateByCategory.Execute()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rows := make([]returnRateByCategoryRowResponse, 0, len(report))
	for _, row := range report {
		rows = append(rows, returnRateByCategoryRowResponse{
			Category:        row.Category,
			ShippedQuantity: row.ShippedQuantity,
			ReturnQuantity:  row.ReturnQuantity,
			ReturnRate:      row.ReturnRate,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(returnRateByCategoryResponse{Categories: rows})
}

func (h ReportHandler) topDiscountedProductsReport(w http.ResponseWriter, r *http.Request) {
	report, err := h.topDiscountedProducts.Execute()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rows := make([]topDiscountedProductRowResponse, 0, len(report))
	for _, row := range report {
		rows = append(rows, topDiscountedProductRowResponse{
			SKU:                 row.SKU,
			ProductName:         row.ProductName,
			QuotedQuantity:      row.QuotedQuantity,
			TotalDiscountAmount: row.TotalDiscountAmount,
			AverageDiscountRate: row.AverageDiscountRate,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(topDiscountedProductsResponse{Products: rows})
}

func (h ReportHandler) lowStockItemsReport(w http.ResponseWriter, r *http.Request) {
	report, err := h.lowStockItems.Execute()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rows := make([]lowStockItemRowResponse, 0, len(report))
	for _, row := range report {
		rows = append(rows, lowStockItemRowResponse{
			SKU:              row.SKU,
			Available:        row.Available,
			ReorderThreshold: row.ReorderThreshold,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(lowStockItemsResponse{Items: rows})
}
