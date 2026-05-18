package http

import (
	"encoding/json"
	"net/http"

	"hexagonal-architecture/internal/core/application"
)

type QuoteHandler struct {
	createQuote application.CreateDraftQuoteUseCase
}

type createQuoteRequest struct {
	CustomerID string `json:"customerId"`
}

type quoteResponse struct {
	ID         string `json:"id"`
	CustomerID string `json:"customerId"`
	Status     string `json:"status"`
}

func NewQuoteHandler(createQuote application.CreateDraftQuoteUseCase) QuoteHandler {
	return QuoteHandler{createQuote: createQuote}
}

func (h QuoteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost || r.URL.Path != "/quotes" {
		http.NotFound(w, r)
		return
	}

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
