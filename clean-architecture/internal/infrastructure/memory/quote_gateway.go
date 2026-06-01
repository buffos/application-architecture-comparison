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
