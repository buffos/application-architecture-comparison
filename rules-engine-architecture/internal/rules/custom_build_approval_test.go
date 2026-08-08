package rules

import (
	"testing"

	"rules-engine-architecture/internal/engine"
)

func TestCustomBuildApprovalRuleActivatesForCustomBuildLine(t *testing.T) {
	memory := engine.NewWorkingMemory(
		engine.CustomerFact{ID: "CUST-001"},
		engine.QuoteFact{
			ID: "Q-1001",
			Lines: []engine.QuoteLineFact{
				{ProductID: "PRD-002", Quantity: 1},
			},
		},
		[]engine.ProductFact{
			{ID: "PRD-002", Category: "CustomBuild"},
		},
	)
	rule := CustomBuildApprovalRule{}

	if !rule.Evaluate(memory) {
		t.Fatal("expected a CustomBuild line to require approval")
	}
	if err := rule.Execute(memory); err != nil {
		t.Fatalf("execute rule: %v", err)
	}
	if len(memory.Findings) != 1 {
		t.Fatalf("expected one finding, got %d", len(memory.Findings))
	}
}
