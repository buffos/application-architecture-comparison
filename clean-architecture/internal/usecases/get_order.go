package usecases

import "clean-architecture/internal/entities"

type GetOrderInput struct {
	OrderID string
}

type GetOrderOutput struct {
	OrderID       string
	CustomerID    string
	SourceQuoteID string
	Status        string
	Lines         int
}

type GetOrderInputBoundary interface {
	Execute(input GetOrderInput) error
}

type GetOrderOutputBoundary interface {
	Present(output GetOrderOutput) error
}

type OrderReader interface {
	FindByID(id string) (entities.Order, error)
}

type GetOrderInteractor struct {
	orders OrderReader
	output GetOrderOutputBoundary
}

func NewGetOrderInteractor(orders OrderReader, output GetOrderOutputBoundary) GetOrderInteractor {
	return GetOrderInteractor{
		orders: orders,
		output: output,
	}
}

func (uc GetOrderInteractor) Execute(input GetOrderInput) error {
	order, err := uc.orders.FindByID(input.OrderID)
	if err != nil {
		return err
	}

	return uc.output.Present(GetOrderOutput{
		OrderID:       order.ID,
		CustomerID:    order.CustomerID,
		SourceQuoteID: order.SourceQuoteID,
		Status:        order.Status,
		Lines:         len(order.Lines),
	})
}
