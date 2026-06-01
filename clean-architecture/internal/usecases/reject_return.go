package usecases

type RejectReturnInput struct {
	ReturnRequestID string
	IdempotencyKey  string
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
	idempotency IdempotencyStore
	returns     ReturnRequestEditor
	output      RejectReturnOutputBoundary
}

func NewRejectReturnInteractor(idempotency IdempotencyStore, returns ReturnRequestEditor, output RejectReturnOutputBoundary) RejectReturnInteractor {
	return RejectReturnInteractor{
		idempotency: idempotency,
		returns:     returns,
		output:      output,
	}
}

func (uc RejectReturnInteractor) Execute(input RejectReturnInput) error {
	if input.IdempotencyKey == "" {
		return ErrIdempotencyKeyRequired
	}

	if existingID, found, err := uc.idempotency.Find("reject-return", input.IdempotencyKey); err != nil {
		return err
	} else if found {
		request, err := uc.returns.FindByID(existingID)
		if err != nil {
			return err
		}

		return uc.output.Present(RejectReturnOutput{
			ReturnRequestID: request.ID,
			OrderID:         request.OrderID,
			Status:          request.Status,
		})
	}

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

	if err := uc.idempotency.Save("reject-return", input.IdempotencyKey, request.ID); err != nil {
		return err
	}

	return uc.output.Present(RejectReturnOutput{
		ReturnRequestID: request.ID,
		OrderID:         request.OrderID,
		Status:          request.Status,
	})
}
