package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"hexagonal-architecture/internal/core/application"
	"hexagonal-architecture/internal/core/domain"
)

type CustomerHandler struct {
	getCustomer  application.GetCustomerUseCase
	listCustomer application.ListCustomersUseCase
}

type customerResponse struct {
	ID     string `json:"id"`
	Active bool   `json:"active"`
}

func NewCustomerHandler(getCustomer application.GetCustomerUseCase, listCustomer application.ListCustomersUseCase) CustomerHandler {
	return CustomerHandler{
		getCustomer:  getCustomer,
		listCustomer: listCustomer,
	}
}

func (h CustomerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet && r.URL.Path == "/customers":
		h.listCustomers(w, r)
	case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/customers/"):
		h.getCustomerRequest(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h CustomerHandler) getCustomerRequest(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/customers/")

	customer, err := h.getCustomer.Execute(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(toCustomerResponse(customer))
}

func (h CustomerHandler) listCustomers(w http.ResponseWriter, r *http.Request) {
	activeOnly := strings.EqualFold(r.URL.Query().Get("status"), "Active")
	customers, err := h.listCustomer.Execute(activeOnly)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := make([]customerResponse, 0, len(customers))
	for _, customer := range customers {
		response = append(response, toCustomerResponse(customer))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func toCustomerResponse(customer domain.Customer) customerResponse {
	return customerResponse{
		ID:     customer.ID,
		Active: customer.Active,
	}
}
