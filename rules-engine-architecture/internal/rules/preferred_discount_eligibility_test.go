package rules

import (
	"testing"

	"rules-engine-architecture/internal/engine"
)

func TestPreferredDiscountEligibilityRulePublishesFactForPreferredCustomer(t *testing.T) {
	memory := engine.NewWorkingMemory(
		engine.CustomerFact{ID: "CUST-001", Tier: "Preferred"},
		engine.QuoteFact{ID: "Q-1001"},
		nil,
	)
	rule := PreferredDiscountEligibilityRule{}

	if !rule.Evaluate(memory) {
		t.Fatal("expected Preferred customer to be eligible")
	}
	if err := rule.Execute(memory); err != nil {
		t.Fatalf("execute rule: %v", err)
	}
	if !memory.HasDerivedFact(engine.PreferredDiscountEligibleFact) {
		t.Fatal("expected the eligibility fact to be published")
	}
	if len(memory.Findings) != 0 {
		t.Fatalf("expected eligibility to add no blocking finding, got %d", len(memory.Findings))
	}
}

func TestPreferredDiscountEligibilityRuleIgnoresOtherCustomerTiers(t *testing.T) {
	memory := engine.NewWorkingMemory(
		engine.CustomerFact{ID: "CUST-002", Tier: "Standard"},
		engine.QuoteFact{ID: "Q-1002"},
		nil,
	)

	if (PreferredDiscountEligibilityRule{}).Evaluate(memory) {
		t.Fatal("expected a non-Preferred customer to be ineligible")
	}
}
