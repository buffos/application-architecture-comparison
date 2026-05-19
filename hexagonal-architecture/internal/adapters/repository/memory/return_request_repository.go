package memory

import (
	"sync"

	"hexagonal-architecture/internal/core/domain"
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

func (r *ReturnRequestRepository) FindByID(id string) (domain.ReturnRequest, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	request, ok := r.requests[id]
	if !ok {
		return domain.ReturnRequest{}, domain.ErrReturnRequestNotFound
	}

	return request, nil
}

func (r *ReturnRequestRepository) ListByStatus(status string) ([]domain.ReturnRequest, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	requests := make([]domain.ReturnRequest, 0, len(r.requests))
	for _, request := range r.requests {
		if status == "" || request.Status == status {
			requests = append(requests, request)
		}
	}

	return requests, nil
}
