package memory

import (
	"sort"
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

func (g *ReturnRequestGateway) ListByStatus(status string) ([]entities.ReturnRequest, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	requests := make([]entities.ReturnRequest, 0, len(g.requests))
	for _, request := range g.requests {
		if status != "" && request.Status != status {
			continue
		}

		requests = append(requests, request)
	}

	sort.Slice(requests, func(i int, j int) bool {
		return requests[i].ID < requests[j].ID
	})

	return requests, nil
}
