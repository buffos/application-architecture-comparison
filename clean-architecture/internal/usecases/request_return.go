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

type RequestReturnInteractor struct {
	orders   OrderEditor
	returns  ReturnRequestWriter
	output   RequestReturnOutputBoundary
}

func NewRequestReturnInteractor(orders OrderEditor, returns ReturnRequestWriter, output RequestReturnOutputBoundary) RequestReturnInteractor {
	return RequestReturnInteractor{
		orders:  orders,
		returns: returns,
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

	if err := uc.returns.Save(request); err != nil {
		return err
	}

	return uc.output.Present(RequestReturnOutput{
		ReturnRequestID: request.ID,
		OrderID:         request.OrderID,
		Status:          request.Status,
	})
}
