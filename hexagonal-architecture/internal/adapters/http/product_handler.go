package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"hexagonal-architecture/internal/core/application"
	"hexagonal-architecture/internal/core/domain"
)

type ProductHandler struct {
	getProduct  application.GetProductUseCase
	listProduct application.ListProductsUseCase
}

type productResponse struct {
	SKU              string `json:"sku"`
	Name             string `json:"name"`
	Category         string `json:"category"`
	Available        bool   `json:"available"`
	ReturnWindowDays int    `json:"returnWindowDays"`
}

func NewProductHandler(getProduct application.GetProductUseCase, listProduct application.ListProductsUseCase) ProductHandler {
	return ProductHandler{
		getProduct:  getProduct,
		listProduct: listProduct,
	}
}

func (h ProductHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet && r.URL.Path == "/products":
		h.listProducts(w, r)
	case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/products/"):
		h.getProductRequest(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h ProductHandler) getProductRequest(w http.ResponseWriter, r *http.Request) {
	sku := strings.TrimPrefix(r.URL.Path, "/products/")

	product, err := h.getProduct.Execute(sku)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(toProductResponse(product))
}

func (h ProductHandler) listProducts(w http.ResponseWriter, r *http.Request) {
	availability := strings.EqualFold(r.URL.Query().Get("availability"), "Available")
	products, err := h.listProduct.Execute(r.URL.Query().Get("category"), availability)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := make([]productResponse, 0, len(products))
	for _, product := range products {
		response = append(response, toProductResponse(product))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func toProductResponse(product domain.Product) productResponse {
	return productResponse{
		SKU:              product.SKU,
		Name:             product.Name,
		Category:         product.Category,
		Available:        product.Available,
		ReturnWindowDays: product.ReturnWindowDays,
	}
}
