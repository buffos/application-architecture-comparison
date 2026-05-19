package application

import (
	"hexagonal-architecture/internal/core/domain"
	"hexagonal-architecture/internal/core/ports"
)

type AcceptReturnUseCase struct {
	returns ports.ReturnRequestRepository
	policy  ports.ReturnEligibilityPolicy
}

func NewAcceptReturnUseCase(returns ports.ReturnRequestRepository, policy ports.ReturnEligibilityPolicy) AcceptReturnUseCase {
	return AcceptReturnUseCase{
		returns: returns,
		policy:  policy,
	}
}

func (uc AcceptReturnUseCase) Execute(returnRequestID string) (domain.ReturnRequest, error) {
	request, err := uc.returns.FindByID(returnRequestID)
	if err != nil {
		return domain.ReturnRequest{}, err
	}

	canAccept, err := uc.policy.CanAccept(request)
	if err != nil {
		return domain.ReturnRequest{}, err
	}

	if !canAccept {
		return domain.ReturnRequest{}, domain.ErrReturnNotEligible
	}

	if err := request.Accept(); err != nil {
		return domain.ReturnRequest{}, err
	}

	if err := uc.returns.Save(request); err != nil {
		return domain.ReturnRequest{}, err
	}

	return request, nil
}
