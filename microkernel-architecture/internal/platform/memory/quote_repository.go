package memory

import (
	"slices"
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

func (r *QuoteRepository) ListByStatus(status string) ([]quotes.Quote, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	results := make([]quotes.Quote, 0)
	for _, quote := range r.quotes {
		if status == "" || quote.Status == status {
			results = append(results, quote)
		}
	}

	slices.SortFunc(results, func(a quotes.Quote, b quotes.Quote) int {
		if a.ID < b.ID {
			return -1
		}
		if a.ID > b.ID {
			return 1
		}
		return 0
	})

	return results, nil
}
