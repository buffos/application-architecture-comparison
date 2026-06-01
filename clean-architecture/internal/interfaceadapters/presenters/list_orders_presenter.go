package presenters

import (
	"fmt"

	"clean-architecture/internal/usecases"
)

type OrderListItemViewModel struct {
	OrderID       string
	CustomerID    string
	SourceQuoteID string
	Status        string
}

type ListOrdersViewModel struct {
	Message string
	Status  string
	Count   int
	Orders  []OrderListItemViewModel
}

type ListOrdersPresenter struct {
	viewModel ListOrdersViewModel
}

func NewListOrdersPresenter() *ListOrdersPresenter {
	return &ListOrdersPresenter{}
}

func (p *ListOrdersPresenter) Present(output usecases.ListOrdersOutput) error {
	items := make([]OrderListItemViewModel, 0, len(output.Orders))
	for _, order := range output.Orders {
		items = append(items, OrderListItemViewModel{
			OrderID:       order.OrderID,
			CustomerID:    order.CustomerID,
			SourceQuoteID: order.SourceQuoteID,
			Status:        order.Status,
		})
	}

	p.viewModel = ListOrdersViewModel{
		Message: fmt.Sprintf("listed orders: status=%s count=%d", output.Status, output.Count),
		Status:  output.Status,
		Count:   output.Count,
		Orders:  items,
	}

	return nil
}

func (p *ListOrdersPresenter) ViewModel() ListOrdersViewModel {
	return p.viewModel
}
