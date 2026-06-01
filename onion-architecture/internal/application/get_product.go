package application

import "onion-architecture/internal/domain"

type ProductFinder interface {
	FindBySKU(sku string) (domain.Product, error)
	List(category string, activeOnly bool) ([]domain.Product, error)
}

type GetProductQuery struct {
	SKU string
}

type ProductDetails struct {
	SKU              string
	Name             string
	Category         string
	Active           bool
	UnitPrice        int
	ReturnWindowDays int
}

type GetProductService struct {
	products ProductFinder
}

func NewGetProductService(products ProductFinder) GetProductService {
	return GetProductService{
		products: products,
	}
}

func (s GetProductService) Execute(query GetProductQuery) (ProductDetails, error) {
	product, err := s.products.FindBySKU(query.SKU)
	if err != nil {
		return ProductDetails{}, err
	}

	return ProductDetails{
		SKU:              product.SKU,
		Name:             product.Name,
		Category:         product.Category,
		Active:           product.Active,
		UnitPrice:        product.UnitPrice,
		ReturnWindowDays: product.ReturnWindowDays,
	}, nil
}
