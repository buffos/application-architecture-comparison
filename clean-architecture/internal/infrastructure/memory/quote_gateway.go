package memory

import (
	"sort"
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

func (g *QuoteGateway) ListByStatus(status string) ([]entities.Quote, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	quotes := make([]entities.Quote, 0, len(g.quotes))
	for _, quote := range g.quotes {
		if status != "" && quote.Status != status {
			continue
		}

		quotes = append(quotes, quote)
	}

	sort.Slice(quotes, func(i int, j int) bool {
		return quotes[i].ID < quotes[j].ID
	})

	return quotes, nil
}
