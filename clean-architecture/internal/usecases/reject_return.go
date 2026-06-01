package usecases

type RejectReturnInput struct {
	ReturnRequestID string
	ReviewedBy      string
	ReviewNote      string
}

type RejectReturnOutput struct {
	ReturnRequestID string
	OrderID         string
	Status          string
}

type RejectReturnInputBoundary interface {
	Execute(input RejectReturnInput) error
}

type RejectReturnOutputBoundary interface {
	Present(output RejectReturnOutput) error
}

type RejectReturnInteractor struct {
	returns ReturnRequestEditor
	output  RejectReturnOutputBoundary
}

func NewRejectReturnInteractor(returns ReturnRequestEditor, output RejectReturnOutputBoundary) RejectReturnInteractor {
	return RejectReturnInteractor{
		returns: returns,
		output:  output,
	}
}

func (uc RejectReturnInteractor) Execute(input RejectReturnInput) error {
	request, err := uc.returns.FindByID(input.ReturnRequestID)
	if err != nil {
		return err
	}

	if err := request.Reject(input.ReviewedBy, input.ReviewNote); err != nil {
		return err
	}

	if err := uc.returns.Save(request); err != nil {
		return err
	}

	return uc.output.Present(RejectReturnOutput{
		ReturnRequestID: request.ID,
		OrderID:         request.OrderID,
		Status:          request.Status,
	})
}
