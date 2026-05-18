package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"hexagonal-architecture/internal/adapters/repository/memory"
	"hexagonal-architecture/internal/core/application"
)

func TestQuoteHandlerCreatesDraftQuote(t *testing.T) {
	repo := memory.NewQuoteRepository()
	createQuote := application.NewCreateDraftQuoteUseCase(repo)
	getQuote := application.NewGetQuoteUseCase(repo)
	handler := NewQuoteHandler(createQuote, getQuote)

	request := httptest.NewRequest(http.MethodPost, "/quotes", strings.NewReader(`{"customerId":"customer-123"}`))
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, recorder.Code)
	}

	body := recorder.Body.String()
	if !strings.Contains(body, `"customerId":"customer-123"`) {
		t.Fatalf("expected response to contain customer id, got %s", body)
	}

	if !strings.Contains(body, `"status":"Draft"`) {
		t.Fatalf("expected response to contain draft status, got %s", body)
	}

	getRequest := httptest.NewRequest(http.MethodGet, "/quotes/quote-001", nil)
	getRecorder := httptest.NewRecorder()

	handler.ServeHTTP(getRecorder, getRequest)

	if getRecorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, getRecorder.Code)
	}

	getBody := getRecorder.Body.String()
	if !strings.Contains(getBody, `"id":"quote-001"`) {
		t.Fatalf("expected fetched response to contain quote id, got %s", getBody)
	}
}
