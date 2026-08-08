package rules

import (
	"testing"

	"rules-engine-architecture/internal/engine"
)

func TestInventoryShortageRuleBackordersWhenProductPolicyAllowsIt(t *testing.T) {
	memory := inventoryMemory(engine.StockShortageBackorder)
	rule := InventoryShortageRule{}

	if !rule.Evaluate(memory) {
		t.Fatal("expected an insufficient-stock line to match")
	}
	if err := rule.Execute(memory); err != nil {
		t.Fatalf("execute rule: %v", err)
	}

	decision := engine.DecisionFromFindings(memory.Findings)
	if decision.FulfillmentAction != engine.FulfillmentBackorder {
		t.Fatalf("expected backorder action, got %s", decision.FulfillmentAction)
	}
	if decision.Outcome != engine.OutcomeAllowed {
		t.Fatalf("expected backorder to remain policy-allowed, got %s", decision.Outcome)
	}
}

func TestInventoryShortageRuleRejectsWhenProductPolicyRequiresIt(t *testing.T) {
	memory := inventoryMemory(engine.StockShortageReject)
	rule := InventoryShortageRule{}

	if err := rule.Execute(memory); err != nil {
		t.Fatalf("execute rule: %v", err)
	}

	decision := engine.DecisionFromFindings(memory.Findings)
	if decision.FulfillmentAction != engine.FulfillmentReject {
		t.Fatalf("expected reject action, got %s", decision.FulfillmentAction)
	}
	if decision.Outcome != engine.OutcomeRejected {
		t.Fatalf("expected rejected outcome, got %s", decision.Outcome)
	}
}

func TestInventoryShortageRuleIgnoresSufficientStock(t *testing.T) {
	memory := inventoryMemory(engine.StockShortageBackorder)
	memory.Products[0].AvailableQuantity = 5

	rule := InventoryShortageRule{}
	if rule.Evaluate(memory) {
		t.Fatal("expected sufficient stock not to match")
	}
}

func inventoryMemory(policy engine.StockShortagePolicy) *engine.WorkingMemory {
	return engine.NewWorkingMemory(
		engine.CustomerFact{ID: "CUST-001"},
		engine.QuoteFact{
			ID: "Q-1001",
			Lines: []engine.QuoteLineFact{
				{ProductID: "PRD-002", Quantity: 3},
			},
		},
		[]engine.ProductFact{
			{
				ID:                "PRD-002",
				AvailableQuantity: 2,
				ShortagePolicy:    policy,
			},
		},
	)
}
