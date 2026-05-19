package application

import (
	"hexagonal-architecture/internal/core/domain"
	"hexagonal-architecture/internal/core/ports"
)

type RejectReturnUseCase struct {
	returns ports.ReturnRequestRepository
}

func NewRejectReturnUseCase(returns ports.ReturnRequestRepository) RejectReturnUseCase {
	return RejectReturnUseCase{returns: returns}
}

func (uc RejectReturnUseCase) Execute(returnRequestID, reviewedBy, reviewNote string) (domain.ReturnRequest, error) {
	request, err := uc.returns.FindByID(returnRequestID)
	if err != nil {
		return domain.ReturnRequest{}, err
	}

	if err := request.Reject(reviewedBy, reviewNote); err != nil {
		return domain.ReturnRequest{}, err
	}

	if err := uc.returns.Save(request); err != nil {
		return domain.ReturnRequest{}, err
	}

	return request, nil
}
