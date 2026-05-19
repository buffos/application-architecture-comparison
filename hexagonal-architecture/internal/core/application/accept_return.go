package application

import (
	"hexagonal-architecture/internal/core/domain"
	"hexagonal-architecture/internal/core/ports"
)

type AcceptReturnUseCase struct {
	returns ports.ReturnRequestRepository
	policy  ports.ReturnEligibilityPolicy
	keys    ports.IdempotencyStore
}

func NewAcceptReturnUseCase(returns ports.ReturnRequestRepository, policy ports.ReturnEligibilityPolicy, keys ports.IdempotencyStore) AcceptReturnUseCase {
	return AcceptReturnUseCase{
		returns: returns,
		policy:  policy,
		keys:    keys,
	}
}

func (uc AcceptReturnUseCase) Execute(returnRequestID, reviewedBy, idempotencyKey string) (domain.ReturnRequest, error) {
	seen, err := uc.keys.Seen("accept-return", idempotencyKey)
	if err != nil {
		return domain.ReturnRequest{}, err
	}
	if seen {
		storedID, err := uc.keys.ResourceID("accept-return", idempotencyKey)
		if err != nil {
			return domain.ReturnRequest{}, err
		}
		return uc.returns.FindByID(storedID)
	}

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

	if err := request.Accept(reviewedBy); err != nil {
		return domain.ReturnRequest{}, err
	}

	if err := uc.returns.Save(request); err != nil {
		return domain.ReturnRequest{}, err
	}

	if err := uc.keys.Remember("accept-return", idempotencyKey, request.ID); err != nil {
		return domain.ReturnRequest{}, err
	}

	return request, nil
}
