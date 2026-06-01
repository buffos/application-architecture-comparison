package presenters

import (
	"fmt"

	"clean-architecture/internal/usecases"
)

type CustomerListItemViewModel struct {
	CustomerID string
	Active     bool
}

type ListCustomersViewModel struct {
	Message    string
	ActiveOnly bool
	Count      int
	Customers  []CustomerListItemViewModel
}

type ListCustomersPresenter struct {
	viewModel ListCustomersViewModel
}

func NewListCustomersPresenter() *ListCustomersPresenter {
	return &ListCustomersPresenter{}
}

func (p *ListCustomersPresenter) Present(output usecases.ListCustomersOutput) error {
	items := make([]CustomerListItemViewModel, 0, len(output.Customers))
	for _, customer := range output.Customers {
		items = append(items, CustomerListItemViewModel{
			CustomerID: customer.CustomerID,
			Active:     customer.Active,
		})
	}

	p.viewModel = ListCustomersViewModel{
		Message:    fmt.Sprintf("listed customers: activeOnly=%t count=%d", output.ActiveOnly, output.Count),
		ActiveOnly: output.ActiveOnly,
		Count:      output.Count,
		Customers:  items,
	}

	return nil
}

func (p *ListCustomersPresenter) ViewModel() ListCustomersViewModel {
	return p.viewModel
}
