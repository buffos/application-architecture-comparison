package products

import "errors"

var ErrProductNotFound = errors.New("product not found")
var ErrProductInactive = errors.New("product is inactive")

type Product struct {
	SKU              string
	Name             string
	Category         string
	Active           bool
	UnitPrice        int
	ReturnWindowDays int
}

type ProductForQuote struct {
	SKU              string
	Name             string
	Category         string
	UnitPrice        int
	ReturnWindowDays int
}

type Repository interface {
	FindBySKU(sku string) (Product, error)
	Save(product Product) error
}
