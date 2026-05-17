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
	customerRepo := memory.NewCustomerRepository()
	productRepo := memory.NewProductRepository()
	customerService := application.NewCustomerService(customerRepo)
	service := application.NewQuoteService(repo, customerRepo, productRepo)
	handler := NewQuoteHandler(service)

	customer, err := customerService.CreateCustomer("Acme Corp", "Preferred", "Invoice30")
	if err != nil {
		t.Fatalf("expected customer creation to succeed, got %v", err)
	}

	createRequest := httptest.NewRequest(http.MethodPost, "/quotes", strings.NewReader(`{"customerId":"`+customer.ID+`"}`))
	createRecorder := httptest.NewRecorder()

	handler.ServeHTTP(createRecorder, createRequest)

	if createRecorder.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, createRecorder.Code)
	}

	body := createRecorder.Body.String()
	if !strings.Contains(body, `"customerId":"`+customer.ID+`"`) {
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
