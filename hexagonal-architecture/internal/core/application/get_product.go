package application

import (
	"hexagonal-architecture/internal/core/domain"
	"hexagonal-architecture/internal/core/ports"
)

type GetProductUseCase struct {
	products ports.ProductLookup
}

func NewGetProductUseCase(products ports.ProductLookup) GetProductUseCase {
	return GetProductUseCase{products: products}
}

func (uc GetProductUseCase) Execute(sku string) (domain.Product, error) {
	return uc.products.FindBySKU(sku)
}
