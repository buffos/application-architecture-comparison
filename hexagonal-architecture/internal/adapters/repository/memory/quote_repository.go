package memory

import (
	"sync"

	"hexagonal-architecture/internal/core/domain"
)

type QuoteRepository struct {
	mu     sync.RWMutex
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

func (r *QuoteRepository) FindByID(id string) (domain.Quote, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	quote, ok := r.quotes[id]
	if !ok {
		return domain.Quote{}, domain.ErrQuoteNotFound
	}

	return quote, nil
}
