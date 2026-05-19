package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"hexagonal-architecture/internal/core/application"
	"hexagonal-architecture/internal/core/domain"
)

type ReturnHandler struct {
	getReturn  application.GetReturnRequestUseCase
	listReturn application.ListReturnRequestsUseCase
}

type returnResponse struct {
	ID          string `json:"id"`
	OrderID     string `json:"orderId"`
	Status      string `json:"status"`
	RequestedBy string `json:"requestedBy"`
	ReviewedBy  string `json:"reviewedBy,omitempty"`
	ProcessedBy string `json:"processedBy,omitempty"`
}

func NewReturnHandler(getReturn application.GetReturnRequestUseCase, listReturn application.ListReturnRequestsUseCase) ReturnHandler {
	return ReturnHandler{
		getReturn:  getReturn,
		listReturn: listReturn,
	}
}

func (h ReturnHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet && r.URL.Path == "/returns":
		h.listReturnRequests(w, r)
	case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/returns/"):
		h.getReturnRequest(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h ReturnHandler) getReturnRequest(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/returns/")

	request, err := h.getReturn.Execute(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(toReturnResponse(request))
}

func (h ReturnHandler) listReturnRequests(w http.ResponseWriter, r *http.Request) {
	requests, err := h.listReturn.Execute(r.URL.Query().Get("status"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := make([]returnResponse, 0, len(requests))
	for _, request := range requests {
		response = append(response, toReturnResponse(request))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func toReturnResponse(request domain.ReturnRequest) returnResponse {
	return returnResponse{
		ID:          request.ID,
		OrderID:     request.OrderID,
		Status:      request.Status,
		RequestedBy: request.RequestedBy,
		ReviewedBy:  request.ReviewedBy,
		ProcessedBy: request.ProcessedBy,
	}
}
