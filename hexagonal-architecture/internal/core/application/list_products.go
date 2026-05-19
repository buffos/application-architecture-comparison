package application

import (
	"hexagonal-architecture/internal/core/domain"
	"hexagonal-architecture/internal/core/ports"
)

type ListProductsUseCase struct {
	products ports.ProductLookup
}

func NewListProductsUseCase(products ports.ProductLookup) ListProductsUseCase {
	return ListProductsUseCase{products: products}
}

func (uc ListProductsUseCase) Execute(category string, availableOnly bool) ([]domain.Product, error) {
	return uc.products.List(category, availableOnly)
}
