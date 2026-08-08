package rules

import (
	"testing"

	"rules-engine-architecture/internal/engine"
)

func TestCustomBuildApprovalRulePublishesApprovalFactForCustomBuildLine(t *testing.T) {
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
	if len(memory.Findings) != 0 {
		t.Fatalf("expected the producer to add no decision finding, got %d", len(memory.Findings))
	}
	if !memory.HasDerivedFact(engine.ManagerApprovalRequiredFact) {
		t.Fatal("expected the producer to publish the manager approval fact")
	}
}
