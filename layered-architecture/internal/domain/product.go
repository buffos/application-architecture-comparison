package domain

import "errors"

var ErrProductSKURequired = errors.New("product sku is required")
var ErrProductNameRequired = errors.New("product name is required")
var ErrProductCategoryRequired = errors.New("product category is required")
var ErrProductNotFound = errors.New("product not found")
var ErrProductUnavailable = errors.New("product is unavailable")

type Product struct {
	ID        string
	SKU       string
	Name      string
	Category  string
	Available bool
}

func NewProduct(sku string, name string, category string, available bool) (Product, error) {
	if sku == "" {
		return Product{}, ErrProductSKURequired
	}

	if name == "" {
		return Product{}, ErrProductNameRequired
	}

	if category == "" {
		return Product{}, ErrProductCategoryRequired
	}

	return Product{
		ID:        sku,
		SKU:       sku,
		Name:      name,
		Category:  category,
		Available: available,
	}, nil
}
