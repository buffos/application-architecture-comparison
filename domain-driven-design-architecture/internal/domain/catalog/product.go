package catalog

import (
	"errors"
	"strings"
)

var (
	ErrProductSKURequired     = errors.New("product sku is required")
	ErrProductNameRequired    = errors.New("product name is required")
	ErrProductCategoryInvalid = errors.New("product category is invalid")
	ErrPriceAmountNegative    = errors.New("price amount cannot be negative")
	ErrPriceCurrencyRequired  = errors.New("price currency is required")
	ErrReturnWindowInvalid    = errors.New("return window must be positive")
	ErrProductAlreadyActive   = errors.New("product is already active")
	ErrProductAlreadyInactive = errors.New("product is already inactive")
	ErrProductInactive        = errors.New("product is inactive")
)

type ProductSKU string
type ProductCategory string

const (
	ProductCategoryStandard    ProductCategory = "Standard"
	ProductCategoryCustomBuild ProductCategory = "CustomBuild"
	ProductCategoryClearance   ProductCategory = "Clearance"
)

// Price is a Catalog-context value object. It is intentionally not shared with
// Quoting's Money type; translation between bounded contexts is explicit.
type Price struct {
	cents    int64
	currency string
}

func NewPrice(cents int64, currency string) (Price, error) {
	if cents < 0 {
		return Price{}, ErrPriceAmountNegative
	}
	currency = strings.ToUpper(strings.TrimSpace(currency))
	if currency == "" {
		return Price{}, ErrPriceCurrencyRequired
	}
	return Price{cents: cents, currency: currency}, nil
}

func (p Price) Cents() int64     { return p.cents }
func (p Price) Currency() string { return p.currency }

// Product is the aggregate root for the Catalog bounded context.
type Product struct {
	sku              ProductSKU
	name             string
	category         ProductCategory
	basePrice        Price
	returnWindowDays int
	active           bool
}

func NewProduct(sku ProductSKU, name string, category ProductCategory, basePrice Price, returnWindowDays int) (Product, error) {
	if sku == "" {
		return Product{}, ErrProductSKURequired
	}
	if strings.TrimSpace(name) == "" {
		return Product{}, ErrProductNameRequired
	}
	if !validCategory(category) {
		return Product{}, ErrProductCategoryInvalid
	}
	if basePrice.Currency() == "" {
		return Product{}, ErrPriceCurrencyRequired
	}
	if returnWindowDays <= 0 {
		return Product{}, ErrReturnWindowInvalid
	}
	return Product{sku: sku, name: name, category: category, basePrice: basePrice, returnWindowDays: returnWindowDays, active: true}, nil
}

func (p Product) SKU() ProductSKU           { return p.sku }
func (p Product) Name() string              { return p.name }
func (p Product) Category() ProductCategory { return p.category }
func (p Product) BasePrice() Price          { return p.basePrice }
func (p Product) ReturnWindowDays() int     { return p.returnWindowDays }
func (p Product) Active() bool              { return p.active }

func (p Product) EnsureSellable() error {
	if !p.active {
		return ErrProductInactive
	}
	return nil
}

func (p *Product) ChangeBasePrice(price Price) error {
	if price.Currency() == "" {
		return ErrPriceCurrencyRequired
	}
	p.basePrice = price
	return nil
}

func (p *Product) Deactivate() error {
	if !p.active {
		return ErrProductAlreadyInactive
	}
	p.active = false
	return nil
}

func (p *Product) Activate() error {
	if p.active {
		return ErrProductAlreadyActive
	}
	p.active = true
	return nil
}

func validCategory(category ProductCategory) bool {
	switch category {
	case ProductCategoryStandard, ProductCategoryCustomBuild, ProductCategoryClearance:
		return true
	default:
		return false
	}
}
