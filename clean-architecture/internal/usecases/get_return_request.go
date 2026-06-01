package usecases

import "clean-architecture/internal/entities"

type GetReturnRequestInput struct {
	ReturnRequestID string
}

type GetReturnRequestOutput struct {
	ReturnRequestID string
	OrderID         string
	Status          string
	Reason          string
	RequestedBy     string
	ReviewedBy      string
	ProcessedBy     string
}

type GetReturnRequestInputBoundary interface {
	Execute(input GetReturnRequestInput) error
}

type GetReturnRequestOutputBoundary interface {
	Present(output GetReturnRequestOutput) error
}

type ReturnRequestReader interface {
	FindByID(id string) (entities.ReturnRequest, error)
}

type GetReturnRequestInteractor struct {
	returns ReturnRequestReader
	output  GetReturnRequestOutputBoundary
}

func NewGetReturnRequestInteractor(returns ReturnRequestReader, output GetReturnRequestOutputBoundary) GetReturnRequestInteractor {
	return GetReturnRequestInteractor{
		returns: returns,
		output:  output,
	}
}

func (uc GetReturnRequestInteractor) Execute(input GetReturnRequestInput) error {
	request, err := uc.returns.FindByID(input.ReturnRequestID)
	if err != nil {
		return err
	}

	return uc.output.Present(GetReturnRequestOutput{
		ReturnRequestID: request.ID,
		OrderID:         request.OrderID,
		Status:          request.Status,
		Reason:          request.Reason,
		RequestedBy:     request.RequestedBy,
		ReviewedBy:      request.ReviewedBy,
		ProcessedBy:     request.ProcessedBy,
	})
}
