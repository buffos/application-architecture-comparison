package ports

import "hexagonal-architecture/internal/core/domain"

type QuoteRepository interface {
	Save(quote domain.Quote) error
	FindByID(id string) (domain.Quote, error)
}
