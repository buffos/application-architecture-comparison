package catalog

import (
	"errors"
	"strings"
)

var (
	ErrProductSKURequired         = errors.New("product sku is required")
	ErrProductNameRequired        = errors.New("product name is required")
	ErrProductCategoryInvalid     = errors.New("product category is invalid")
	ErrProductInactive            = errors.New("product is not sellable")
	ErrProductAlreadyDiscontinued = errors.New("product is already discontinued")
)

type SKU string
type ProductCategory string

const (
	ProductCategoryStandard    ProductCategory = "Standard"
	ProductCategoryCustomBuild ProductCategory = "CustomBuild"
	ProductCategoryClearance   ProductCategory = "Clearance"
)

// Product is a rich Catalog domain object. Its state is private so catalog
// rules stay behind behavior rather than becoming public field assignments.
type Product struct {
	sku       SKU
	name      string
	category  ProductCategory
	basePrice Price
	active    bool
}

func NewProduct(sku SKU, name string, category ProductCategory, basePrice Price) (Product, error) {
	if sku == "" {
		return Product{}, ErrProductSKURequired
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return Product{}, ErrProductNameRequired
	}
	if !validProductCategory(category) {
		return Product{}, ErrProductCategoryInvalid
	}
	if basePrice.Currency() == "" {
		return Product{}, ErrPriceCurrencyRequired
	}

	return Product{
		sku:       sku,
		name:      name,
		category:  category,
		basePrice: basePrice,
		active:    true,
	}, nil
}

func (product Product) SKU() SKU                  { return product.sku }
func (product Product) Name() string              { return product.name }
func (product Product) Category() ProductCategory { return product.category }
func (product Product) BasePrice() Price          { return product.basePrice }
func (product Product) Active() bool              { return product.active }

func (product Product) EnsureSellable() error {
	if !product.active {
		return ErrProductInactive
	}
	return nil
}

func (product *Product) ChangePrice(price Price) error {
	if price.Currency() == "" {
		return ErrPriceCurrencyRequired
	}

	product.basePrice = price
	return nil
}

func (product *Product) Discontinue() error {
	if !product.active {
		return ErrProductAlreadyDiscontinued
	}

	product.active = false
	return nil
}

func validProductCategory(category ProductCategory) bool {
	switch category {
	case ProductCategoryStandard, ProductCategoryCustomBuild, ProductCategoryClearance:
		return true
	default:
		return false
	}
}
