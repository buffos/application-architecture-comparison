package application

import (
	"errors"
	"strings"
	"sync"

	"rules-engine-architecture/internal/engine"
)

// IdempotencyStore remembers completed command results. A production adapter
// can replace this in-memory implementation without changing the Rule Engine.
type IdempotencyStore struct {
	mu        sync.Mutex
	decisions map[string]engine.PolicyDecision
}

func NewIdempotencyStore() *IdempotencyStore {
	return &IdempotencyStore{decisions: map[string]engine.PolicyDecision{}}
}

func (store *IdempotencyStore) Load(key string) (engine.PolicyDecision, bool) {
	store.mu.Lock()
	defer store.mu.Unlock()

	decision, found := store.decisions[key]
	if !found {
		return engine.PolicyDecision{}, false
	}

	return cloneDecision(decision), true
}

func (store *IdempotencyStore) Save(key string, decision engine.PolicyDecision) {
	store.mu.Lock()
	defer store.mu.Unlock()

	store.decisions[key] = cloneDecision(decision)
}

// ReturnDecisionService adds command history around a deterministic Rule
// Engine evaluation.
type ReturnDecisionService struct {
	ruleEngine *engine.Engine
	store      *IdempotencyStore
}

func NewReturnDecisionService(
	ruleEngine *engine.Engine,
	store *IdempotencyStore,
) *ReturnDecisionService {
	return &ReturnDecisionService{
		ruleEngine: ruleEngine,
		store:      store,
	}
}

func (service *ReturnDecisionService) Evaluate(
	commandKey string,
	memory *engine.WorkingMemory,
	maxCycles int,
) (engine.PolicyDecision, int, bool, error) {
	if strings.TrimSpace(commandKey) == "" {
		return engine.PolicyDecision{}, 0, false, errors.New("return command key is required")
	}

	if decision, found := service.store.Load(commandKey); found {
		return decision, 0, true, nil
	}

	decision, cycles, err := service.ruleEngine.DecideUntilStable(memory, maxCycles)
	if err != nil {
		return engine.PolicyDecision{}, cycles, false, err
	}

	service.store.Save(commandKey, decision)
	return decision, cycles, false, nil
}

func cloneDecision(decision engine.PolicyDecision) engine.PolicyDecision {
	decision.Reasons = append([]engine.Finding(nil), decision.Reasons...)
	decision.RequiredReviews = append(
		[]engine.ReviewRequirement(nil),
		decision.RequiredReviews...,
	)
	return decision
}
