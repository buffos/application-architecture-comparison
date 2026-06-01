package memory

import (
	"sync"

	"clean-architecture/internal/entities"
)

type ReturnRequestGateway struct {
	mu       sync.RWMutex
	requests map[string]entities.ReturnRequest
}

func NewReturnRequestGateway() *ReturnRequestGateway {
	return &ReturnRequestGateway{
		requests: make(map[string]entities.ReturnRequest),
	}
}

func (g *ReturnRequestGateway) Save(request entities.ReturnRequest) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.requests[request.ID] = request
	return nil
}

func (g *ReturnRequestGateway) FindByID(id string) (entities.ReturnRequest, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	request, ok := g.requests[id]
	if !ok {
		return entities.ReturnRequest{}, entities.ErrQuoteNotFound
	}

	return request, nil
}
