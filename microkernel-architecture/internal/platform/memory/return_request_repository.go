package memory

import (
	"sync"

	"microkernel-architecture/internal/plugins/returns"
)

type ReturnRequestRepository struct {
	mu       sync.RWMutex
	requests map[string]returns.ReturnRequest
}

func NewReturnRequestRepository() *ReturnRequestRepository {
	return &ReturnRequestRepository{
		requests: make(map[string]returns.ReturnRequest),
	}
}

func (r *ReturnRequestRepository) Save(request returns.ReturnRequest) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.requests[request.ID] = request
	return nil
}

func (r *ReturnRequestRepository) FindByID(id string) (returns.ReturnRequest, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	request, ok := r.requests[id]
	if !ok {
		return returns.ReturnRequest{}, returns.ErrReturnRequestNotFound
	}

	return request, nil
}
