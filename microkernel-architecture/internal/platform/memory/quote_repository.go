package memory

import (
	"sync"

	"microkernel-architecture/internal/plugins/quotes"
)

type QuoteRepository struct {
	mu     sync.RWMutex
	quotes map[string]quotes.Quote
}

func NewQuoteRepository() *QuoteRepository {
	return &QuoteRepository{
		quotes: make(map[string]quotes.Quote),
	}
}

func (r *QuoteRepository) Save(quote quotes.Quote) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.quotes[quote.ID] = quote
	return nil
}

func (r *QuoteRepository) FindByID(id string) (quotes.Quote, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	quote, ok := r.quotes[id]
	if !ok {
		return quotes.Quote{}, quotes.ErrQuoteNotFound
	}

	return quote, nil
}
