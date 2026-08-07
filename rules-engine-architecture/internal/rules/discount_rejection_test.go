package rules

import "testing"

func TestDiscountRejectionRuleActivatesAboveMaximum(t *testing.T) {
	memory := newMemoryWithDiscount(30)
	rule := DiscountRejectionRule{}

	if !rule.Evaluate(memory) {
		t.Fatal("expected a 30% discount to be rejected")
	}
	if err := rule.Execute(memory); err != nil {
		t.Fatalf("execute rule: %v", err)
	}
	if len(memory.Findings) != 1 {
		t.Fatalf("expected one finding, got %d", len(memory.Findings))
	}
}
