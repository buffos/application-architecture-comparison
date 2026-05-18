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
	handler := NewQuoteHandler(createQuote)

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
}
