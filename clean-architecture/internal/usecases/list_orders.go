package usecases

import "clean-architecture/internal/entities"

type ListOrdersInput struct {
	Status string
}

type OrderListItem struct {
	OrderID       string
	CustomerID    string
	SourceQuoteID string
	Status        string
}

type ListOrdersOutput struct {
	Status string
	Count  int
	Orders []OrderListItem
}

type ListOrdersInputBoundary interface {
	Execute(input ListOrdersInput) error
}

type ListOrdersOutputBoundary interface {
	Present(output ListOrdersOutput) error
}

type OrderLister interface {
	ListByStatus(status string) ([]entities.Order, error)
}

type ListOrdersInteractor struct {
	orders OrderLister
	output ListOrdersOutputBoundary
}

func NewListOrdersInteractor(orders OrderLister, output ListOrdersOutputBoundary) ListOrdersInteractor {
	return ListOrdersInteractor{
		orders: orders,
		output: output,
	}
}

func (uc ListOrdersInteractor) Execute(input ListOrdersInput) error {
	orders, err := uc.orders.ListByStatus(input.Status)
	if err != nil {
		return err
	}

	items := make([]OrderListItem, 0, len(orders))
	for _, order := range orders {
		items = append(items, OrderListItem{
			OrderID:       order.ID,
			CustomerID:    order.CustomerID,
			SourceQuoteID: order.SourceQuoteID,
			Status:        order.Status,
		})
	}

	return uc.output.Present(ListOrdersOutput{
		Status: input.Status,
		Count:  len(items),
		Orders: items,
	})
}
