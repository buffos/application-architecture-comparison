package usecases

import "clean-architecture/internal/entities"

type AcceptReturnInput struct {
	ReturnRequestID string
}

type AcceptReturnOutput struct {
	ReturnRequestID string
	OrderID         string
	Status          string
}

type AcceptReturnInputBoundary interface {
	Execute(input AcceptReturnInput) error
}

type AcceptReturnOutputBoundary interface {
	Present(output AcceptReturnOutput) error
}

type ReturnRequestEditor interface {
	FindByID(id string) (entities.ReturnRequest, error)
	Save(request entities.ReturnRequest) error
}

type RefundGateway interface {
	Refund(order entities.Order) error
}

type InventoryRestock interface {
	Restock(items []entities.InventoryReservationItem) error
}

type AcceptReturnInteractor struct {
	orders   OrderEditor
	returns  ReturnRequestEditor
	policy   ReturnEligibilityPolicy
	refunds  RefundGateway
	restock  InventoryRestock
	output   AcceptReturnOutputBoundary
}

func NewAcceptReturnInteractor(orders OrderEditor, returns ReturnRequestEditor, policy ReturnEligibilityPolicy, refunds RefundGateway, restock InventoryRestock, output AcceptReturnOutputBoundary) AcceptReturnInteractor {
	return AcceptReturnInteractor{
		orders:  orders,
		returns: returns,
		policy:  policy,
		refunds: refunds,
		restock: restock,
		output:  output,
	}
}

func (uc AcceptReturnInteractor) Execute(input AcceptReturnInput) error {
	request, err := uc.returns.FindByID(input.ReturnRequestID)
	if err != nil {
		return err
	}

	order, err := uc.orders.FindByID(request.OrderID)
	if err != nil {
		return err
	}

	allowed, err := uc.policy.CanAccept(order, request)
	if err != nil {
		return err
	}
	if !allowed {
		return entities.ErrQuoteCannotTransition
	}

	if err := request.Accept(); err != nil {
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

	if err := request.MarkRefunded(); err != nil {
		return err
	}

	if err := uc.returns.Save(request); err != nil {
		return err
	}

	return uc.output.Present(AcceptReturnOutput{
		ReturnRequestID: request.ID,
		OrderID:         request.OrderID,
		Status:          request.Status,
	})
}
