package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"layered-architecture/internal/application"
)

type QuoteHandler struct {
	service application.QuoteService
}

type createQuoteRequest struct {
	CustomerID string `json:"customerId"`
}

type quoteResponse struct {
	ID         string `json:"id"`
	CustomerID string `json:"customerId"`
	Status     string `json:"status"`
	LineCount  int    `json:"lineCount"`
}

func NewQuoteHandler(service application.QuoteService) QuoteHandler {
	return QuoteHandler{service: service}
}

func (h QuoteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && r.URL.Path == "/quotes":
		h.createQuote(w, r)
	case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/quotes/"):
		h.getQuote(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h QuoteHandler) createQuote(w http.ResponseWriter, r *http.Request) {
	var req createQuoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	quote, err := h.service.CreateDraftQuote(req.CustomerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeQuoteJSON(w, quoteResponse{
		ID:         quote.ID,
		CustomerID: quote.CustomerID,
		Status:     quote.Status,
		LineCount:  len(quote.Lines),
	}, http.StatusCreated)
}

func (h QuoteHandler) getQuote(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/quotes/")

	quote, err := h.service.GetQuote(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	writeQuoteJSON(w, quoteResponse{
		ID:         quote.ID,
		CustomerID: quote.CustomerID,
		Status:     quote.Status,
		LineCount:  len(quote.Lines),
	}, http.StatusOK)
}

func writeQuoteJSON(w http.ResponseWriter, response quoteResponse, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(response)
}
