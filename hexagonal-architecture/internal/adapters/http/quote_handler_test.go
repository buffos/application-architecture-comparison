package http

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"hexagonal-architecture/internal/adapters/repository/memory"
	"hexagonal-architecture/internal/adapters/services/approval"
	"hexagonal-architecture/internal/adapters/services/pricing"
	"hexagonal-architecture/internal/core/application"
	"hexagonal-architecture/internal/core/domain"
)

func TestQuoteHandlerCreatesDraftQuote(t *testing.T) {
	quoteRepo := memory.NewQuoteRepository()
	customerRepo := memory.NewCustomerRepository()
	productRepo := memory.NewProductRepository()
	pricingPolicy := pricing.NewFixedPricingPolicy()
	approvalPolicy := approval.NewCategoryApprovalPolicy()
	if err := customerRepo.Save(domain.Customer{ID: "customer-123", Active: true}); err != nil {
		t.Fatalf("expected customer save to succeed, got %v", err)
	}
	if err := productRepo.Save(domain.Product{SKU: "DESK-001", Name: "Executive Desk", Category: "CustomBuild", BasePrice: 50000, Available: true, ReturnWindowDays: 30}); err != nil {
		t.Fatalf("expected product save to succeed, got %v", err)
	}

	createQuote := application.NewCreateDraftQuoteUseCase(quoteRepo, customerRepo)
	getQuote := application.NewGetQuoteUseCase(quoteRepo)
	listQuotes := application.NewListQuotesUseCase(quoteRepo)
	addQuoteLine := application.NewAddQuoteLineUseCase(quoteRepo, productRepo, pricingPolicy)
	submitQuote := application.NewSubmitQuoteUseCase(quoteRepo, approvalPolicy)
	handler := NewQuoteHandler(createQuote, getQuote, listQuotes)

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

	customQuote, err := createQuote.Execute("customer-123")
	if err != nil {
		t.Fatalf("expected second quote creation to succeed, got %v", err)
	}

	if _, err := addQuoteLine.Execute(customQuote.ID, "DESK-001", 1); err != nil {
		t.Fatalf("expected add quote line to succeed, got %v", err)
	}

	if _, err := submitQuote.Execute(customQuote.ID); err != nil {
		t.Fatalf("expected quote submission to succeed, got %v", err)
	}

	listRequest := httptest.NewRequest(http.MethodGet, "/quotes?status=PendingApproval", nil)
	listRecorder := httptest.NewRecorder()
	handler.ServeHTTP(listRecorder, listRequest)

	if listRecorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, listRecorder.Code)
	}

	listBody := listRecorder.Body.String()
	if !strings.Contains(listBody, `"id":"`+customQuote.ID+`"`) || !strings.Contains(listBody, `"status":"PendingApproval"`) {
		t.Fatalf("expected list response to contain pending approval quote, got %s", listBody)
	}
}
