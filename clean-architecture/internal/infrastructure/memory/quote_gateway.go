package memory

import (
	"sync"

	"clean-architecture/internal/entities"
)

type QuoteGateway struct {
	mu     sync.RWMutex
	quotes map[string]entities.Quote
}

func NewQuoteGateway() *QuoteGateway {
	return &QuoteGateway{
		quotes: make(map[string]entities.Quote),
	}
}

func (g *QuoteGateway) Save(quote entities.Quote) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.quotes[quote.ID] = quote
	return nil
}

func (g *QuoteGateway) FindByID(id string) (entities.Quote, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	quote, ok := g.quotes[id]
	if !ok {
		return entities.Quote{}, entities.ErrQuoteNotFound
	}

	return quote, nil
}
