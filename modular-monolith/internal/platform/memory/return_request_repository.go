package memory

import (
	"sync"

	"modular-monolith/internal/modules/returns"
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

func (r *ReturnRequestRepository) ListByStatus(status string) ([]returns.ReturnRequest, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]returns.ReturnRequest, 0, len(r.requests))
	for _, request := range r.requests {
		if status == "" || request.Status == status {
			list = append(list, request)
		}
	}

	return list, nil
}
