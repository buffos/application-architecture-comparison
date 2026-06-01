package usecases

import "clean-architecture/internal/entities"

type GetProductInput struct {
	SKU string
}

type GetProductOutput struct {
	SKU              string
	Name             string
	Category         string
	BasePrice        int
	Available        bool
	ReturnWindowDays int
}

type GetProductInputBoundary interface {
	Execute(input GetProductInput) error
}

type GetProductOutputBoundary interface {
	Present(output GetProductOutput) error
}

type ProductReader interface {
	FindBySKU(sku string) (entities.Product, error)
}

type GetProductInteractor struct {
	products ProductReader
	output   GetProductOutputBoundary
}

func NewGetProductInteractor(products ProductReader, output GetProductOutputBoundary) GetProductInteractor {
	return GetProductInteractor{
		products: products,
		output:   output,
	}
}

func (uc GetProductInteractor) Execute(input GetProductInput) error {
	product, err := uc.products.FindBySKU(input.SKU)
	if err != nil {
		return err
	}

	return uc.output.Present(GetProductOutput{
		SKU:              product.SKU,
		Name:             product.Name,
		Category:         product.Category,
		BasePrice:        product.BasePrice,
		Available:        product.Available,
		ReturnWindowDays: product.ReturnWindowDays,
	})
}
