package presenters

import (
	"fmt"

	"clean-architecture/internal/usecases"
)

type GetCustomerViewModel struct {
	Message    string
	CustomerID string
	Active     bool
}

type GetCustomerPresenter struct {
	viewModel GetCustomerViewModel
}

func NewGetCustomerPresenter() *GetCustomerPresenter {
	return &GetCustomerPresenter{}
}

func (p *GetCustomerPresenter) Present(output usecases.GetCustomerOutput) error {
	p.viewModel = GetCustomerViewModel{
		Message:    fmt.Sprintf("loaded customer: id=%s active=%t", output.CustomerID, output.Active),
		CustomerID: output.CustomerID,
		Active:     output.Active,
	}

	return nil
}

func (p *GetCustomerPresenter) ViewModel() GetCustomerViewModel {
	return p.viewModel
}
