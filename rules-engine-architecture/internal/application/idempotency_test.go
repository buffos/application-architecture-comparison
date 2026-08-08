package application

import (
	"testing"

	"rules-engine-architecture/internal/engine"
	"rules-engine-architecture/internal/rules"
)

func TestReturnDecisionServiceReplaysStoredDecision(t *testing.T) {
	ruleEngine := engine.NewEngine()
	ruleEngine.Register(rules.ReturnPolicyRule{})
	service := NewReturnDecisionService(ruleEngine, NewIdempotencyStore())

	firstMemory := returnDecisionMemory("Standard")
	first, cycles, replayed, err := service.Evaluate("return-001", firstMemory, 5)
	if err != nil {
		t.Fatalf("first evaluation: %v", err)
	}
	if replayed {
		t.Fatal("expected the first evaluation not to be replayed")
	}
	if cycles == 0 {
		t.Fatal("expected the first evaluation to execute the Rule Engine")
	}

	retryMemory := returnDecisionMemory("Clearance")
	second, retryCycles, replayed, err := service.Evaluate("return-001", retryMemory, 5)
	if err != nil {
		t.Fatalf("retry evaluation: %v", err)
	}
	if !replayed {
		t.Fatal("expected the second evaluation to be replayed")
	}
	if retryCycles != 0 {
		t.Fatalf("expected no inference cycles on replay, got %d", retryCycles)
	}
	if second.ReturnAction != first.ReturnAction {
		t.Fatalf("expected the stored action %s, got %s", first.ReturnAction, second.ReturnAction)
	}
	if len(retryMemory.Trace) != 0 {
		t.Fatalf("expected the replay to skip the Rule Engine, got %d traces", len(retryMemory.Trace))
	}
}

func TestReturnDecisionServiceRequiresCommandKey(t *testing.T) {
	ruleEngine := engine.NewEngine()
	service := NewReturnDecisionService(ruleEngine, NewIdempotencyStore())

	_, _, _, err := service.Evaluate(" ", returnDecisionMemory("Standard"), 5)
	if err == nil {
		t.Fatal("expected a missing command key to fail")
	}
}

func returnDecisionMemory(category string) *engine.WorkingMemory {
	memory := engine.NewWorkingMemory(
		engine.CustomerFact{ID: "CUST-001"},
		engine.QuoteFact{ID: "Q-1001"},
		[]engine.ProductFact{{ID: "PRD-002", Category: category}},
	)
	memory.Order = engine.OrderFact{ID: "ORD-1001", Status: engine.OrderShipped}
	memory.ReturnRequest = engine.ReturnRequestFact{
		Requested:         true,
		ProductID:         "PRD-002",
		Quantity:          1,
		ShippedQuantity:   1,
		DaysSinceShipment: 5,
		ReturnWindowDays:  30,
		RequestedBy: engine.ActorFact{
			ID:   "warehouse-clerk-001",
			Role: "warehouse-clerk",
		},
	}
	return memory
}
