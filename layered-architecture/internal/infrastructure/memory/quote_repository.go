package memory

import (
	"sync"

	"layered-architecture/internal/domain"
)

type QuoteRepository struct {
	mu     sync.Mutex
	quotes map[string]domain.Quote
}

func NewQuoteRepository() *QuoteRepository {
	return &QuoteRepository{
		quotes: make(map[string]domain.Quote),
	}
}

func (r *QuoteRepository) Save(quote domain.Quote) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.quotes[quote.ID] = quote

	return nil
}
