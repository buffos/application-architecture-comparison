package usecases

import "clean-architecture/internal/entities"

type ListProductsInput struct {
	Category      string
	AvailableOnly bool
}

type ProductListItem struct {
	SKU              string
	Name             string
	Category         string
	BasePrice        int
	Available        bool
	ReturnWindowDays int
}

type ListProductsOutput struct {
	Category      string
	AvailableOnly bool
	Count         int
	Products      []ProductListItem
}

type ListProductsInputBoundary interface {
	Execute(input ListProductsInput) error
}

type ListProductsOutputBoundary interface {
	Present(output ListProductsOutput) error
}

type ProductLister interface {
	List(category string, availableOnly bool) ([]entities.Product, error)
}

type ListProductsInteractor struct {
	products ProductLister
	output   ListProductsOutputBoundary
}

func NewListProductsInteractor(products ProductLister, output ListProductsOutputBoundary) ListProductsInteractor {
	return ListProductsInteractor{
		products: products,
		output:   output,
	}
}

func (uc ListProductsInteractor) Execute(input ListProductsInput) error {
	products, err := uc.products.List(input.Category, input.AvailableOnly)
	if err != nil {
		return err
	}

	items := make([]ProductListItem, 0, len(products))
	for _, product := range products {
		items = append(items, ProductListItem{
			SKU:              product.SKU,
			Name:             product.Name,
			Category:         product.Category,
			BasePrice:        product.BasePrice,
			Available:        product.Available,
			ReturnWindowDays: product.ReturnWindowDays,
		})
	}

	return uc.output.Present(ListProductsOutput{
		Category:      input.Category,
		AvailableOnly: input.AvailableOnly,
		Count:         len(items),
		Products:      items,
	})
}
