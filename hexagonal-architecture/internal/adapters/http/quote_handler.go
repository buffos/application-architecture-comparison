package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"hexagonal-architecture/internal/core/application"
)

type QuoteHandler struct {
	createQuote application.CreateDraftQuoteUseCase
	getQuote    application.GetQuoteUseCase
}

type createQuoteRequest struct {
	CustomerID string `json:"customerId"`
}

type quoteResponse struct {
	ID         string `json:"id"`
	CustomerID string `json:"customerId"`
	Status     string `json:"status"`
}

func NewQuoteHandler(createQuote application.CreateDraftQuoteUseCase, getQuote application.GetQuoteUseCase) QuoteHandler {
	return QuoteHandler{
		createQuote: createQuote,
		getQuote:    getQuote,
	}
}

func (h QuoteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && r.URL.Path == "/quotes":
		h.createQuoteRequest(w, r)
	case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/quotes/"):
		h.getQuoteRequest(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h QuoteHandler) createQuoteRequest(w http.ResponseWriter, r *http.Request) {
	var req createQuoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	quote, err := h.createQuote.Execute(req.CustomerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(quoteResponse{
		ID:         quote.ID,
		CustomerID: quote.CustomerID,
		Status:     quote.Status,
	})
}

func (h QuoteHandler) getQuoteRequest(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/quotes/")

	quote, err := h.getQuote.Execute(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(quoteResponse{
		ID:         quote.ID,
		CustomerID: quote.CustomerID,
		Status:     quote.Status,
	})
}
