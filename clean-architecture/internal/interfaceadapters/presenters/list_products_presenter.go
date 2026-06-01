package presenters

import (
	"fmt"

	"clean-architecture/internal/usecases"
)

type ProductListItemViewModel struct {
	SKU              string
	Name             string
	Category         string
	BasePrice        int
	Available        bool
	ReturnWindowDays int
}

type ListProductsViewModel struct {
	Message       string
	Category      string
	AvailableOnly bool
	Count         int
	Products      []ProductListItemViewModel
}

type ListProductsPresenter struct {
	viewModel ListProductsViewModel
}

func NewListProductsPresenter() *ListProductsPresenter {
	return &ListProductsPresenter{}
}

func (p *ListProductsPresenter) Present(output usecases.ListProductsOutput) error {
	items := make([]ProductListItemViewModel, 0, len(output.Products))
	for _, product := range output.Products {
		items = append(items, ProductListItemViewModel{
			SKU:              product.SKU,
			Name:             product.Name,
			Category:         product.Category,
			BasePrice:        product.BasePrice,
			Available:        product.Available,
			ReturnWindowDays: product.ReturnWindowDays,
		})
	}

	p.viewModel = ListProductsViewModel{
		Message:       fmt.Sprintf("listed products: category=%s availableOnly=%t count=%d", output.Category, output.AvailableOnly, output.Count),
		Category:      output.Category,
		AvailableOnly: output.AvailableOnly,
		Count:         output.Count,
		Products:      items,
	}

	return nil
}

func (p *ListProductsPresenter) ViewModel() ListProductsViewModel {
	return p.viewModel
}
