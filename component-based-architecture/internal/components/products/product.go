package products

import "errors"

var (
	ErrProductSKURequired = errors.New("product sku is required")
	ErrProductNotFound    = errors.New("product not found")
	ErrProductInactive    = errors.New("product is inactive")
)

type Product struct {
	SKU       string
	Name      string
	Category  string
	Active    bool
	UnitPrice int
}

// ProductForQuote is the product snapshot that crosses the component
// boundary. It contains only the data Quotes needs for this workflow.
type ProductForQuote struct {
	SKU       string
	Name      string
	Category  string
	UnitPrice int
}
