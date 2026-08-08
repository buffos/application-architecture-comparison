package plugins

import (
	"testing"

	"rules-engine-architecture/internal/engine"
)

func TestSeasonalSurchargeRuleContributesThroughRuleContract(t *testing.T) {
	memory := engine.NewWorkingMemory(
		engine.CustomerFact{ID: "CUST-001"},
		engine.QuoteFact{
			ID:    "Q-1001",
			Lines: []engine.QuoteLineFact{{ProductID: "PRD-002", Quantity: 2, UnitPriceCents: 12500}},
		},
		[]engine.ProductFact{{ID: "PRD-002", Category: "CustomBuild"}},
	)
	rule := NewSeasonalSurchargeRule("CustomBuild", 5)

	if !rule.Evaluate(memory) {
		t.Fatal("expected the plugin Rule to match the configured category")
	}
	if err := rule.Execute(memory); err != nil {
		t.Fatalf("execute plugin Rule: %v", err)
	}

	decision := engine.DecisionFromFindings(memory.Findings)
	if decision.PricingAdjustmentCents != 1250 {
		t.Fatalf("expected 1250 cents surcharge, got %d", decision.PricingAdjustmentCents)
	}
}

func TestSeasonalSurchargeRuleIgnoresOtherCategories(t *testing.T) {
	memory := engine.NewWorkingMemory(
		engine.CustomerFact{ID: "CUST-001"},
		engine.QuoteFact{
			ID:    "Q-1001",
			Lines: []engine.QuoteLineFact{{ProductID: "PRD-001", Quantity: 1, UnitPriceCents: 12500}},
		},
		[]engine.ProductFact{{ID: "PRD-001", Category: "Standard"}},
	)

	if (NewSeasonalSurchargeRule("CustomBuild", 5)).Evaluate(memory) {
		t.Fatal("expected the plugin Rule not to match Standard products")
	}
}
