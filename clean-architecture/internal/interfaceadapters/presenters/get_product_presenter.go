package presenters

import (
	"fmt"

	"clean-architecture/internal/usecases"
)

type GetProductViewModel struct {
	Message          string
	SKU              string
	Name             string
	Category         string
	BasePrice        int
	Available        bool
	ReturnWindowDays int
}

type GetProductPresenter struct {
	viewModel GetProductViewModel
}

func NewGetProductPresenter() *GetProductPresenter {
	return &GetProductPresenter{}
}

func (p *GetProductPresenter) Present(output usecases.GetProductOutput) error {
	p.viewModel = GetProductViewModel{
		Message:          fmt.Sprintf("loaded product: sku=%s category=%s available=%t", output.SKU, output.Category, output.Available),
		SKU:              output.SKU,
		Name:             output.Name,
		Category:         output.Category,
		BasePrice:        output.BasePrice,
		Available:        output.Available,
		ReturnWindowDays: output.ReturnWindowDays,
	}

	return nil
}

func (p *GetProductPresenter) ViewModel() GetProductViewModel {
	return p.viewModel
}
