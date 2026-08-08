package rules

import (
	"testing"

	"rules-engine-architecture/internal/engine"
)

func TestHighValuePaymentReviewRuleActivatesAboveThreshold(t *testing.T) {
	memory := engine.NewWorkingMemory(
		engine.CustomerFact{ID: "CUST-001"},
		engine.QuoteFact{
			ID: "Q-1001",
			Lines: []engine.QuoteLineFact{
				{ProductID: "PRD-002", Quantity: 1, UnitPriceCents: 125000},
			},
		},
		nil,
	)
	rule := NewHighValuePaymentReviewRule(100000)

	if !rule.Evaluate(memory) {
		t.Fatal("expected the 1250.00 subtotal to require payment review")
	}
	if err := rule.Execute(memory); err != nil {
		t.Fatalf("execute rule: %v", err)
	}
	if len(memory.Findings) != 1 {
		t.Fatalf("expected one finding, got %d", len(memory.Findings))
	}
}

func TestHighValuePaymentReviewRuleIgnoresSubtotalAtThreshold(t *testing.T) {
	memory := engine.NewWorkingMemory(
		engine.CustomerFact{ID: "CUST-001"},
		engine.QuoteFact{
			ID: "Q-1001",
			Lines: []engine.QuoteLineFact{
				{ProductID: "PRD-001", Quantity: 1, UnitPriceCents: 100000},
			},
		},
		nil,
	)
	rule := NewHighValuePaymentReviewRule(100000)

	if rule.Evaluate(memory) {
		t.Fatal("expected a subtotal at the threshold not to require review")
	}
}
