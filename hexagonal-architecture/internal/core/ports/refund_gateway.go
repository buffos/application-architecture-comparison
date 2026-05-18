package ports

import "hexagonal-architecture/internal/core/domain"

type RefundGateway interface {
	Refund(request domain.ReturnRequest) (bool, error)
}
