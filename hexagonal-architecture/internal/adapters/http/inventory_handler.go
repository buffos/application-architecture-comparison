package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"hexagonal-architecture/internal/core/application"
	"hexagonal-architecture/internal/core/domain"
)

type InventoryHandler struct {
	receiveStock    application.ReceiveStockUseCase
	adjustThreshold application.AdjustReorderThresholdUseCase
	getStockRecord  application.GetStockRecordUseCase
}

type inventoryResponse struct {
	SKU              string `json:"sku"`
	Available        int    `json:"available"`
	ReorderThreshold int    `json:"reorderThreshold"`
}

func NewInventoryHandler(
	receiveStock application.ReceiveStockUseCase,
	adjustThreshold application.AdjustReorderThresholdUseCase,
	getStockRecord application.GetStockRecordUseCase,
) InventoryHandler {
	return InventoryHandler{
		receiveStock:    receiveStock,
		adjustThreshold: adjustThreshold,
		getStockRecord:  getStockRecord,
	}
}

func (h InventoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/inventory/") && !strings.Contains(strings.TrimPrefix(r.URL.Path, "/inventory/"), "/"):
		h.getStock(w, r)
	case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/receive"):
		h.receive(w, r)
	case r.Method == http.MethodPatch && strings.HasSuffix(r.URL.Path, "/reorder-threshold"):
		h.adjustReorderThreshold(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h InventoryHandler) getStock(w http.ResponseWriter, r *http.Request) {
	sku := strings.TrimPrefix(r.URL.Path, "/inventory/")

	record, err := h.getStockRecord.Execute(sku)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	writeInventoryResponse(w, http.StatusOK, record)
}

func (h InventoryHandler) receive(w http.ResponseWriter, r *http.Request) {
	sku := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/inventory/"), "/receive")

	var body struct {
		Quantity int `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	record, err := h.receiveStock.Execute(sku, body.Quantity)
	if err != nil {
		http.Error(w, err.Error(), inventoryStatusCode(err))
		return
	}

	writeInventoryResponse(w, http.StatusOK, record)
}

func (h InventoryHandler) adjustReorderThreshold(w http.ResponseWriter, r *http.Request) {
	sku := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/inventory/"), "/reorder-threshold")

	var body struct {
		ReorderThreshold int `json:"reorderThreshold"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	record, err := h.adjustThreshold.Execute(sku, body.ReorderThreshold)
	if err != nil {
		http.Error(w, err.Error(), inventoryStatusCode(err))
		return
	}

	writeInventoryResponse(w, http.StatusOK, record)
}

func writeInventoryResponse(w http.ResponseWriter, status int, record domain.StockRecord) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(inventoryResponse{
		SKU:              record.SKU,
		Available:        record.Available,
		ReorderThreshold: record.ReorderThreshold,
	})
}

func inventoryStatusCode(err error) int {
	switch err {
	case domain.ErrStockRecordNotFound:
		return http.StatusNotFound
	case domain.ErrStockQuantityInvalid, domain.ErrReorderThresholdInvalid:
		return http.StatusBadRequest
	default:
		return http.StatusBadRequest
	}
}
