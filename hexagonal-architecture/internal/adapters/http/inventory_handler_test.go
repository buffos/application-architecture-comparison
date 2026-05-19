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

func TestInventoryHandlerReceivesAdjustsAndGetsStock(t *testing.T) {
	productRepo := memory.NewProductRepository()
	inventory := memory.NewInventoryReservationAdapter(map[string]int{
		"CHAIR-001": 2,
	})

	_ = productRepo.Save(domain.Product{SKU: "CHAIR-001", Name: "Office Chair", Category: "Standard", BasePrice: 10000, Available: true, ReturnWindowDays: 30})

	handler := NewInventoryHandler(
		application.NewReceiveStockUseCase(productRepo, inventory),
		application.NewAdjustReorderThresholdUseCase(inventory),
		application.NewGetStockRecordUseCase(inventory),
	)

	receiveRequest := httptest.NewRequest(http.MethodPost, "/inventory/CHAIR-001/receive", strings.NewReader(`{"quantity":3}`))
	receiveRecorder := httptest.NewRecorder()
	handler.ServeHTTP(receiveRecorder, receiveRequest)

	if receiveRecorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, receiveRecorder.Code)
	}

	if !strings.Contains(receiveRecorder.Body.String(), `"available":5`) {
		t.Fatalf("expected received stock body to contain available 5, got %s", receiveRecorder.Body.String())
	}

	thresholdRequest := httptest.NewRequest(http.MethodPatch, "/inventory/CHAIR-001/reorder-threshold", strings.NewReader(`{"reorderThreshold":3}`))
	thresholdRecorder := httptest.NewRecorder()
	handler.ServeHTTP(thresholdRecorder, thresholdRequest)

	if thresholdRecorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, thresholdRecorder.Code)
	}

	getRequest := httptest.NewRequest(http.MethodGet, "/inventory/CHAIR-001", nil)
	getRecorder := httptest.NewRecorder()
	handler.ServeHTTP(getRecorder, getRequest)

	if getRecorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, getRecorder.Code)
	}

	body := getRecorder.Body.String()
	if !strings.Contains(body, `"sku":"CHAIR-001"`) || !strings.Contains(body, `"reorderThreshold":3`) {
		t.Fatalf("expected stock body to contain sku and threshold, got %s", body)
	}
}
