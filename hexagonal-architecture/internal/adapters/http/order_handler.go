package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"hexagonal-architecture/internal/core/application"
	"hexagonal-architecture/internal/core/domain"
)

type OrderHandler struct {
	getOrder  application.GetOrderUseCase
	listOrder application.ListOrdersUseCase
}

type orderResponse struct {
	ID            string `json:"id"`
	SourceQuoteID string `json:"sourceQuoteId"`
	CustomerID    string `json:"customerId"`
	Status        string `json:"status"`
	PaymentStatus string `json:"paymentStatus"`
}

func NewOrderHandler(getOrder application.GetOrderUseCase, listOrder application.ListOrdersUseCase) OrderHandler {
	return OrderHandler{
		getOrder:  getOrder,
		listOrder: listOrder,
	}
}

func (h OrderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet && r.URL.Path == "/orders":
		h.listOrders(w, r)
	case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/orders/"):
		h.getOrderRequest(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h OrderHandler) getOrderRequest(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/orders/")

	order, err := h.getOrder.Execute(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(toOrderResponse(order))
}

func (h OrderHandler) listOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.listOrder.Execute(r.URL.Query().Get("status"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := make([]orderResponse, 0, len(orders))
	for _, order := range orders {
		response = append(response, toOrderResponse(order))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func toOrderResponse(order domain.Order) orderResponse {
	return orderResponse{
		ID:            order.ID,
		SourceQuoteID: order.SourceQuoteID,
		CustomerID:    order.CustomerID,
		Status:        order.Status,
		PaymentStatus: order.PaymentStatus,
	}
}
