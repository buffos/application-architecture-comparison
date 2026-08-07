package engine_test

import (
	"testing"

	"rules-engine-architecture/internal/engine"
	"rules-engine-architecture/internal/rules"
)

func TestEngineExecutesRegisteredRules(t *testing.T) {
	memory := engine.NewWorkingMemory(
		engine.CustomerFact{ID: "CUST-001"},
		engine.QuoteFact{
			ID:              "Q-1001",
			CustomerID:      "CUST-001",
			DiscountPercent: 20,
		},
		nil,
	)

	ruleEngine := engine.NewEngine()
	ruleEngine.Register(rules.DiscountApprovalRule{})

	if err := ruleEngine.ExecuteAll(memory); err != nil {
		t.Fatalf("execute rules: %v", err)
	}
	if len(memory.Findings) != 1 {
		t.Fatalf("expected one finding, got %d", len(memory.Findings))
	}
}
