package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"hexagonal-architecture/internal/adapters/repository/memory"
	"hexagonal-architecture/internal/core/application"
	"hexagonal-architecture/internal/core/domain"
)

func TestCustomerHandlerGetsAndListsCustomers(t *testing.T) {
	customerRepo := memory.NewCustomerRepository()
	getCustomer := application.NewGetCustomerUseCase(customerRepo)
	listCustomers := application.NewListCustomersUseCase(customerRepo)
	handler := NewCustomerHandler(getCustomer, listCustomers)

	_ = customerRepo.Save(domain.Customer{ID: "customer-001", Active: true})
	_ = customerRepo.Save(domain.Customer{ID: "customer-002", Active: false})

	getRequest := httptest.NewRequest(http.MethodGet, "/customers/customer-001", nil)
	getRecorder := httptest.NewRecorder()
	handler.ServeHTTP(getRecorder, getRequest)

	if getRecorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, getRecorder.Code)
	}

	if !strings.Contains(getRecorder.Body.String(), `"id":"customer-001"`) {
		t.Fatalf("expected customer body to contain id, got %s", getRecorder.Body.String())
	}

	listRequest := httptest.NewRequest(http.MethodGet, "/customers?status=Active", nil)
	listRecorder := httptest.NewRecorder()
	handler.ServeHTTP(listRecorder, listRequest)

	if listRecorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, listRecorder.Code)
	}

	body := listRecorder.Body.String()
	if !strings.Contains(body, `"id":"customer-001"`) || strings.Contains(body, `"id":"customer-002"`) {
		t.Fatalf("expected only active customer in list, got %s", body)
	}
}
