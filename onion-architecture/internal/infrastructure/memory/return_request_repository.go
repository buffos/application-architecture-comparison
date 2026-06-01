package memory

import (
	"sync"

	"onion-architecture/internal/domain"
)

type ReturnRequestRepository struct {
	mu       sync.RWMutex
	requests map[string]domain.ReturnRequest
}

func NewReturnRequestRepository() *ReturnRequestRepository {
	return &ReturnRequestRepository{
		requests: make(map[string]domain.ReturnRequest),
	}
}

func (r *ReturnRequestRepository) Save(request domain.ReturnRequest) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.requests[request.ID] = request
	return nil
}
