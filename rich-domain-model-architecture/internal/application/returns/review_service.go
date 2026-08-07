package returns

import (
	"errors"

	domainreturns "rich-domain-model-architecture/internal/domain/returns"
)

var ErrIdempotencyKeyRequired = errors.New("idempotency key is required")

type ReviewResult struct {
	ReturnRequestID domainreturns.ReturnRequestID
	Status          domainreturns.ReturnStatus
	ReviewedBy      string
	ProcessedBy     string
}

type IdempotencyStore interface {
	Find(key string) (ReviewResult, bool)
	Save(key string, result ReviewResult)
}

type InMemoryIdempotencyStore struct {
	results map[string]ReviewResult
}

func NewInMemoryIdempotencyStore() *InMemoryIdempotencyStore {
	return &InMemoryIdempotencyStore{results: make(map[string]ReviewResult)}
}

func (store *InMemoryIdempotencyStore) Find(key string) (ReviewResult, bool) {
	result, ok := store.results[key]
	return result, ok
}

func (store *InMemoryIdempotencyStore) Save(key string, result ReviewResult) {
	store.results[key] = result
}

type ReviewService struct {
	store IdempotencyStore
}

func NewReviewService(store IdempotencyStore) ReviewService {
	return ReviewService{store: store}
}

func (service ReviewService) Review(request *domainreturns.ReturnRequest, decision domainreturns.ReviewDecision, reviewer, processor, key string) (ReviewResult, error) {
	if key == "" {
		return ReviewResult{}, ErrIdempotencyKeyRequired
	}
	if result, ok := service.store.Find(key); ok {
		return result, nil
	}
	if err := request.ReviewBy(decision, reviewer); err != nil {
		return ReviewResult{}, err
	}
	if decision == domainreturns.ReviewDecisionAccept {
		if err := request.ProcessBy(processor); err != nil {
			return ReviewResult{}, err
		}
	}
	result := ReviewResult{
		ReturnRequestID: request.ID(),
		Status:          request.Status(),
		ReviewedBy:      request.ReviewedBy(),
		ProcessedBy:     request.ProcessedBy(),
	}
	service.store.Save(key, result)
	return result, nil
}

var _ IdempotencyStore = (*InMemoryIdempotencyStore)(nil)
