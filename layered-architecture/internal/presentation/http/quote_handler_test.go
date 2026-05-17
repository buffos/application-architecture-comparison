package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"layered-architecture/internal/application"
	"layered-architecture/internal/infrastructure/memory"
)

func TestQuoteHandlerCreateAndGetQuote(t *testing.T) {
	repo := memory.NewQuoteRepository()
	service := application.NewQuoteService(repo)
	handler := NewQuoteHandler(service)

	createRequest := httptest.NewRequest(http.MethodPost, "/quotes", strings.NewReader(`{"customerId":"customer-123"}`))
	createRecorder := httptest.NewRecorder()

	handler.ServeHTTP(createRecorder, createRequest)

	if createRecorder.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, createRecorder.Code)
	}

	body := createRecorder.Body.String()
	if !strings.Contains(body, `"customerId":"customer-123"`) {
		t.Fatalf("expected created response to contain customer id, got %s", body)
	}

	if !strings.Contains(body, `"id":"quote-001"`) {
		t.Fatalf("expected created response to contain quote id, got %s", body)
	}

	getRequest := httptest.NewRequest(http.MethodGet, "/quotes/quote-001", nil)
	getRecorder := httptest.NewRecorder()

	handler.ServeHTTP(getRecorder, getRequest)

	if getRecorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, getRecorder.Code)
	}

	getBody := getRecorder.Body.String()
	if !strings.Contains(getBody, `"status":"Draft"`) {
		t.Fatalf("expected fetched response to contain draft status, got %s", getBody)
	}
}
