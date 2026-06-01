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
