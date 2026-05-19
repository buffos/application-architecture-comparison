package http

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"hexagonal-architecture/internal/adapters/repository/memory"
	"hexagonal-architecture/internal/core/application"
	"hexagonal-architecture/internal/core/domain"
)

func TestQuoteHandlerCreatesDraftQuote(t *testing.T) {
	quoteRepo := memory.NewQuoteRepository()
	customerRepo := memory.NewCustomerRepository()
	if err := customerRepo.Save(domain.Customer{ID: "customer-123", Active: true}); err != nil {
		t.Fatalf("expected customer save to succeed, got %v", err)
	}

	createQuote := application.NewCreateDraftQuoteUseCase(quoteRepo, customerRepo)
	getQuote := application.NewGetQuoteUseCase(quoteRepo)
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

	matches := regexp.MustCompile(`"id":"([^"]+)"`).FindStringSubmatch(body)
	if len(matches) != 2 {
		t.Fatalf("expected response to contain quote id, got %s", body)
	}

	getRequest := httptest.NewRequest(http.MethodGet, "/quotes/"+matches[1], nil)
	getRecorder := httptest.NewRecorder()

	handler.ServeHTTP(getRecorder, getRequest)

	if getRecorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, getRecorder.Code)
	}

	getBody := getRecorder.Body.String()
	if !strings.Contains(getBody, `"id":"`+matches[1]+`"`) {
		t.Fatalf("expected fetched response to contain quote id, got %s", getBody)
	}
}
