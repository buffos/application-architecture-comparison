package application

import (
	"hexagonal-architecture/internal/core/domain"
	"hexagonal-architecture/internal/core/ports"
)

type AcceptReturnUseCase struct {
	returns ports.ReturnRequestRepository
}

func NewAcceptReturnUseCase(returns ports.ReturnRequestRepository) AcceptReturnUseCase {
	return AcceptReturnUseCase{returns: returns}
}

func (uc AcceptReturnUseCase) Execute(returnRequestID string) (domain.ReturnRequest, error) {
	request, err := uc.returns.FindByID(returnRequestID)
	if err != nil {
		return domain.ReturnRequest{}, err
	}

	if err := request.Accept(); err != nil {
		return domain.ReturnRequest{}, err
	}

	if err := uc.returns.Save(request); err != nil {
		return domain.ReturnRequest{}, err
	}

	return request, nil
}
