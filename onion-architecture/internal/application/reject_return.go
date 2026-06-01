package application

type RejectReturnCommand struct {
	ReturnRequestID string
	IdempotencyKey  string
	ReviewedBy      string
	ReviewNote      string
}

type RejectReturnResult struct {
	ReturnRequestID string
	Status          string
}

type RejectReturnService struct {
	returns ReturnRequestStore
	idempotency IdempotencyStore
}

func NewRejectReturnService(returns ReturnRequestStore, idempotency IdempotencyStore) RejectReturnService {
	return RejectReturnService{
		returns: returns,
		idempotency: idempotency,
	}
}

func (s RejectReturnService) Execute(command RejectReturnCommand) (RejectReturnResult, error) {
	if status, ok, err := s.idempotency.Get(command.IdempotencyKey); err != nil {
		return RejectReturnResult{}, err
	} else if ok {
		return RejectReturnResult{
			ReturnRequestID: command.ReturnRequestID,
			Status:          status,
		}, nil
	}

	request, err := s.returns.FindByID(command.ReturnRequestID)
	if err != nil {
		return RejectReturnResult{}, err
	}

	if err := request.Reject(command.ReviewedBy, command.ReviewNote); err != nil {
		return RejectReturnResult{}, err
	}

	if err := s.returns.Save(request); err != nil {
		return RejectReturnResult{}, err
	}

	if err := s.idempotency.Save(command.IdempotencyKey, request.Status); err != nil {
		return RejectReturnResult{}, err
	}

	return RejectReturnResult{
		ReturnRequestID: request.ID,
		Status:          request.Status,
	}, nil
}
