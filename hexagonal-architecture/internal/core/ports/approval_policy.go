package ports

import "hexagonal-architecture/internal/core/domain"

type ApprovalPolicy interface {
	RequiresApproval(quote domain.Quote) (bool, error)
}
