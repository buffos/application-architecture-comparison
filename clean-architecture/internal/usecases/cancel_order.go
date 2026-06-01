package usecases

import "clean-architecture/internal/entities"

type CancelOrderInput struct {
	OrderID string
}

type CancelOrderOutput struct {
	OrderID string
	Status  string
	Lines   int
}

type CancelOrderInputBoundary interface {
	Execute(input CancelOrderInput) error
}

type CancelOrderOutputBoundary interface {
	Present(output CancelOrderOutput) error
}

type InventoryRelease interface {
	Release(items []entities.InventoryReservationItem) error
}

type CancelOrderInteractor struct {
	orders    OrderEditor
	inventory InventoryRelease
	output    CancelOrderOutputBoundary
}

func NewCancelOrderInteractor(orders OrderEditor, inventory InventoryRelease, output CancelOrderOutputBoundary) CancelOrderInteractor {
	return CancelOrderInteractor{
		orders:    orders,
		inventory: inventory,
		output:    output,
	}
}

func (uc CancelOrderInteractor) Execute(input CancelOrderInput) error {
	order, err := uc.orders.FindByID(input.OrderID)
	if err != nil {
		return err
	}

	if err := order.Cancel(); err != nil {
		return err
	}

	items := make([]entities.InventoryReservationItem, 0, len(order.Lines))
	for _, line := range order.Lines {
		items = append(items, entities.InventoryReservationItem{
			SKU:      line.SKU,
			Quantity: line.Quantity,
		})
	}

	if err := uc.inventory.Release(items); err != nil {
		return err
	}

	if err := uc.orders.Save(order); err != nil {
		return err
	}

	return uc.output.Present(CancelOrderOutput{
		OrderID: order.ID,
		Status:  order.Status,
		Lines:   len(order.Lines),
	})
}
