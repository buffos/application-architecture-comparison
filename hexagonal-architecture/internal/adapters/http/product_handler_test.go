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

func TestProductHandlerGetsAndListsProducts(t *testing.T) {
	productRepo := memory.NewProductRepository()
	getProduct := application.NewGetProductUseCase(productRepo)
	listProducts := application.NewListProductsUseCase(productRepo)
	handler := NewProductHandler(getProduct, listProducts)

	_ = productRepo.Save(domain.Product{SKU: "CHAIR-001", Name: "Office Chair", Category: "Standard", BasePrice: 10000, Available: true, ReturnWindowDays: 30})
	_ = productRepo.Save(domain.Product{SKU: "LAMP-001", Name: "Clearance Lamp", Category: "Clearance", BasePrice: 4000, Available: false, ReturnWindowDays: 0})

	getRequest := httptest.NewRequest(http.MethodGet, "/products/CHAIR-001", nil)
	getRecorder := httptest.NewRecorder()
	handler.ServeHTTP(getRecorder, getRequest)

	if getRecorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, getRecorder.Code)
	}

	if !strings.Contains(getRecorder.Body.String(), `"sku":"CHAIR-001"`) {
		t.Fatalf("expected product body to contain sku, got %s", getRecorder.Body.String())
	}

	listRequest := httptest.NewRequest(http.MethodGet, "/products?category=Standard&availability=Available", nil)
	listRecorder := httptest.NewRecorder()
	handler.ServeHTTP(listRecorder, listRequest)

	if listRecorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, listRecorder.Code)
	}

	body := listRecorder.Body.String()
	if !strings.Contains(body, `"sku":"CHAIR-001"`) || strings.Contains(body, `"sku":"LAMP-001"`) {
		t.Fatalf("expected only available standard product in list, got %s", body)
	}
}
