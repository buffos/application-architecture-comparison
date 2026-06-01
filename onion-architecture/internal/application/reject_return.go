package application

type RejectReturnCommand struct {
	ReturnRequestID string
}

type RejectReturnResult struct {
	ReturnRequestID string
	Status          string
}

type RejectReturnService struct {
	returns ReturnRequestStore
}

func NewRejectReturnService(returns ReturnRequestStore) RejectReturnService {
	return RejectReturnService{
		returns: returns,
	}
}

func (s RejectReturnService) Execute(command RejectReturnCommand) (RejectReturnResult, error) {
	request, err := s.returns.FindByID(command.ReturnRequestID)
	if err != nil {
		return RejectReturnResult{}, err
	}

	if err := request.Reject(); err != nil {
		return RejectReturnResult{}, err
	}

	if err := s.returns.Save(request); err != nil {
		return RejectReturnResult{}, err
	}

	return RejectReturnResult{
		ReturnRequestID: request.ID,
		Status:          request.Status,
	}, nil
}
