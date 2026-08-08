package rules

import (
	"testing"

	"rules-engine-architecture/internal/engine"
)

func TestApprovalWorkflowGateRuleConsumesDerivedApprovalFact(t *testing.T) {
	memory := engine.NewWorkingMemory(
		engine.CustomerFact{ID: "CUST-001"},
		engine.QuoteFact{ID: "Q-1001"},
		nil,
	)
	rule := ApprovalWorkflowGateRule{}

	if rule.Evaluate(memory) {
		t.Fatal("expected the gate to wait for a derived Fact")
	}
	memory.AddDerivedFact(engine.DerivedFact{
		Name:  engine.ManagerApprovalRequiredFact,
		Value: memory.Quote.ID,
	})
	if !rule.Evaluate(memory) {
		t.Fatal("expected the gate to activate after the approval Fact was derived")
	}
	if err := rule.Execute(memory); err != nil {
		t.Fatalf("execute rule: %v", err)
	}
	if len(memory.Findings) != 1 {
		t.Fatalf("expected one workflow finding, got %d", len(memory.Findings))
	}
}
