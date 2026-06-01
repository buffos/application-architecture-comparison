package usecases

import "clean-architecture/internal/entities"

type RequestReturnInput struct {
	OrderID string
	Reason  string
}

type RequestReturnOutput struct {
	ReturnRequestID string
	OrderID         string
	Status          string
}

type RequestReturnInputBoundary interface {
	Execute(input RequestReturnInput) error
}

type RequestReturnOutputBoundary interface {
	Present(output RequestReturnOutput) error
}

type ReturnRequestWriter interface {
	Save(request entities.ReturnRequest) error
}

type RefundGateway interface {
	Refund(order entities.Order) error
}

type InventoryRestock interface {
	Restock(items []entities.InventoryReservationItem) error
}

type RequestReturnInteractor struct {
	orders   OrderEditor
	returns  ReturnRequestWriter
	refunds  RefundGateway
	restock  InventoryRestock
	output   RequestReturnOutputBoundary
}

func NewRequestReturnInteractor(orders OrderEditor, returns ReturnRequestWriter, refunds RefundGateway, restock InventoryRestock, output RequestReturnOutputBoundary) RequestReturnInteractor {
	return RequestReturnInteractor{
		orders:  orders,
		returns: returns,
		refunds: refunds,
		restock: restock,
		output:  output,
	}
}

func (uc RequestReturnInteractor) Execute(input RequestReturnInput) error {
	order, err := uc.orders.FindByID(input.OrderID)
	if err != nil {
		return err
	}

	request, err := entities.NewReturnRequestFromShippedOrder(order, input.Reason)
	if err != nil {
		return err
	}

	if err := uc.refunds.Refund(order); err != nil {
		return err
	}

	items := make([]entities.InventoryReservationItem, 0, len(order.Lines))
	for _, line := range order.Lines {
		items = append(items, entities.InventoryReservationItem{
			SKU:      line.SKU,
			Quantity: line.Quantity,
		})
	}

	if err := uc.restock.Restock(items); err != nil {
		return err
	}

	if err := uc.returns.Save(request); err != nil {
		return err
	}

	return uc.output.Present(RequestReturnOutput{
		ReturnRequestID: request.ID,
		OrderID:         request.OrderID,
		Status:          request.Status,
	})
}
