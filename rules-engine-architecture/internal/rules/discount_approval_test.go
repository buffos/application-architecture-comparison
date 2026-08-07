package rules

import (
	"testing"

	"rules-engine-architecture/internal/engine"
)

func newMemoryWithDiscount(discountPercent int) *engine.WorkingMemory {
	return engine.NewWorkingMemory(
		engine.CustomerFact{ID: "CUST-001", Tier: "Preferred"},
		engine.QuoteFact{
			ID:              "Q-1001",
			CustomerID:      "CUST-001",
			DiscountPercent: discountPercent,
			Status:          "Draft",
		},
		nil,
	)
}

func TestDiscountApprovalRuleAddsFindingAboveThreshold(t *testing.T) {
	memory := newMemoryWithDiscount(20)
	rule := DiscountApprovalRule{}

	if !rule.Evaluate(memory) {
		t.Fatal("expected a 20% discount to require approval")
	}
	if err := rule.Execute(memory); err != nil {
		t.Fatalf("execute rule: %v", err)
	}
	if len(memory.Findings) != 1 {
		t.Fatalf("expected one finding, got %d", len(memory.Findings))
	}
}

func TestDiscountApprovalRuleIgnoresDiscountAtThreshold(t *testing.T) {
	memory := newMemoryWithDiscount(15)
	rule := DiscountApprovalRule{}

	if rule.Evaluate(memory) {
		t.Fatal("expected a 15% discount not to require approval")
	}
}
