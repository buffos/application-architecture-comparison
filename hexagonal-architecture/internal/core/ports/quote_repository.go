package ports

import "hexagonal-architecture/internal/core/domain"

type QuoteRepository interface {
	Save(quote domain.Quote) error
}
