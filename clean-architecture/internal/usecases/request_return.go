package usecases

import "clean-architecture/internal/entities"

type RequestReturnInput struct {
	OrderID     string
	Reason      string
	Lines       []RequestReturnLineInput
	RequestedBy string
}

type RequestReturnLineInput struct {
	SKU      string
	Quantity int
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
	orders  OrderEditor
	returns ReturnRequestWriter
	clock   Clock
	output  RequestReturnOutputBoundary
}

func NewRequestReturnInteractor(orders OrderEditor, returns ReturnRequestWriter, clock Clock, output RequestReturnOutputBoundary) RequestReturnInteractor {
	return RequestReturnInteractor{
		orders:  orders,
		returns: returns,
		clock:   clock,
		output:  output,
	}
}

func (uc RequestReturnInteractor) Execute(input RequestReturnInput) error {
	order, err := uc.orders.FindByID(input.OrderID)
	if err != nil {
		return err
	}

	requestLines := make([]entities.ReturnRequestLine, 0, len(input.Lines))
	for _, line := range input.Lines {
		productName := ""
		for _, orderLine := range order.Lines {
			if orderLine.SKU == line.SKU {
				productName = orderLine.ProductName
				break
			}
		}

		requestLines = append(requestLines, entities.ReturnRequestLine{
			SKU:         line.SKU,
			ProductName: productName,
			Quantity:    line.Quantity,
		})
	}

	request, err := entities.NewReturnRequestFromShippedOrder(order, input.Reason, requestLines, uc.clock.Now(), input.RequestedBy)
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
