package products

import (
	"errors"
	"sort"

	"rich-domain-model-architecture/internal/domain/catalog"
)

var ErrProductNotFound = errors.New("product not found")

type Reader interface {
	GetProduct(sku string) (ProductDetails, error)
	ListProducts(active *bool) []ProductSummary
}

type ProductDetails struct {
	SKU              string
	Name             string
	Category         string
	BasePriceCents   int64
	Currency         string
	ReturnWindowDays int
	Active           bool
}

type ProductSummary struct {
	SKU      string
	Name     string
	Category string
	Active   bool
}

type InMemoryReader struct {
	products map[string]ProductDetails
}

func NewInMemoryReader() *InMemoryReader {
	return &InMemoryReader{products: make(map[string]ProductDetails)}
}

func (reader *InMemoryReader) Save(product catalog.Product) {
	price := product.BasePrice()
	reader.products[string(product.SKU())] = ProductDetails{
		SKU:              string(product.SKU()),
		Name:             product.Name(),
		Category:         string(product.Category()),
		BasePriceCents:   price.Cents(),
		Currency:         price.Currency(),
		ReturnWindowDays: product.ReturnWindowDays(),
		Active:           product.Active(),
	}
}

func (reader *InMemoryReader) GetProduct(sku string) (ProductDetails, error) {
	product, ok := reader.products[sku]
	if !ok {
		return ProductDetails{}, ErrProductNotFound
	}
	return product, nil
}

func (reader *InMemoryReader) ListProducts(active *bool) []ProductSummary {
	result := make([]ProductSummary, 0, len(reader.products))
	for _, product := range reader.products {
		if active != nil && product.Active != *active {
			continue
		}
		result = append(result, ProductSummary{SKU: product.SKU, Name: product.Name, Category: product.Category, Active: product.Active})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].SKU < result[j].SKU })
	return result
}

var _ Reader = (*InMemoryReader)(nil)
