package domain

import "errors"

var ErrProductNotFound = errors.New("product not found")
var ErrProductUnavailable = errors.New("product is unavailable")

type Product struct {
	SKU       string
	Name      string
	Category  string
	BasePrice int
	Available bool
}
