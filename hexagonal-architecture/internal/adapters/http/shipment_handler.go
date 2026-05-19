package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"hexagonal-architecture/internal/core/application"
	"hexagonal-architecture/internal/core/domain"
)

type ShipmentHandler struct {
	getShipment  application.GetShipmentUseCase
	listShipment application.ListShipmentsUseCase
}

type shipmentResponse struct {
	ID      string `json:"id"`
	OrderID string `json:"orderId"`
	Status  string `json:"status"`
	Lines   int    `json:"lines"`
}

func NewShipmentHandler(getShipment application.GetShipmentUseCase, listShipment application.ListShipmentsUseCase) ShipmentHandler {
	return ShipmentHandler{
		getShipment:  getShipment,
		listShipment: listShipment,
	}
}

func (h ShipmentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet && r.URL.Path == "/shipments":
		h.listShipments(w, r)
	case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/shipments/"):
		h.getShipmentRequest(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h ShipmentHandler) getShipmentRequest(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/shipments/")

	shipment, err := h.getShipment.Execute(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(toShipmentResponse(shipment))
}

func (h ShipmentHandler) listShipments(w http.ResponseWriter, r *http.Request) {
	shipments, err := h.listShipment.Execute(r.URL.Query().Get("orderId"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := make([]shipmentResponse, 0, len(shipments))
	for _, shipment := range shipments {
		response = append(response, toShipmentResponse(shipment))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func toShipmentResponse(shipment domain.Shipment) shipmentResponse {
	return shipmentResponse{
		ID:      shipment.ID,
		OrderID: shipment.OrderID,
		Status:  shipment.Status,
		Lines:   len(shipment.Lines),
	}
}
