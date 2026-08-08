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

func TestEngineUsesHigherPriorityRuleForConflictGroup(t *testing.T) {
	memory := engine.NewWorkingMemory(
		engine.CustomerFact{ID: "CUST-001"},
		engine.QuoteFact{
			ID:              "Q-1001",
			CustomerID:      "CUST-001",
			DiscountPercent: 30,
		},
		nil,
	)

	ruleEngine := engine.NewEngine()
	// Register the lower-priority rule first to prove that registration order
	// does not override the conflict policy.
	ruleEngine.Register(rules.DiscountApprovalRule{})
	ruleEngine.Register(rules.DiscountRejectionRule{})

	if err := ruleEngine.ExecuteAll(memory); err != nil {
		t.Fatalf("execute rules: %v", err)
	}
	if len(memory.Findings) != 1 {
		t.Fatalf("expected one finding, got %d", len(memory.Findings))
	}
	if memory.Findings[0].RuleName != "Discount Rejection Rule" {
		t.Fatalf("expected rejection rule, got %q", memory.Findings[0].RuleName)
	}
}

func TestEngineRunsIndependentRulesTogether(t *testing.T) {
	memory := engine.NewWorkingMemory(
		engine.CustomerFact{ID: "CUST-001"},
		engine.QuoteFact{
			ID:              "Q-1001",
			CustomerID:      "CUST-001",
			DiscountPercent: 20,
			Lines: []engine.QuoteLineFact{
				{ProductID: "PRD-002", Quantity: 1},
			},
		},
		[]engine.ProductFact{
			{ID: "PRD-002", Category: "CustomBuild"},
		},
	)

	ruleEngine := engine.NewEngine()
	ruleEngine.Register(rules.DiscountApprovalRule{})
	ruleEngine.Register(rules.CustomBuildApprovalRule{})

	if err := ruleEngine.ExecuteAll(memory); err != nil {
		t.Fatalf("execute rules: %v", err)
	}
	if len(memory.Findings) != 2 {
		t.Fatalf("expected two independent findings, got %d", len(memory.Findings))
	}
}

func TestEngineCanDisableARegisteredRule(t *testing.T) {
	memory := engine.NewWorkingMemory(
		engine.CustomerFact{ID: "CUST-001"},
		engine.QuoteFact{
			ID:              "Q-1001",
			CustomerID:      "CUST-001",
			DiscountPercent: 20,
			Lines: []engine.QuoteLineFact{
				{ProductID: "PRD-002", Quantity: 1},
			},
		},
		[]engine.ProductFact{
			{ID: "PRD-002", Category: "CustomBuild"},
		},
	)

	ruleEngine := engine.NewEngine()
	ruleEngine.Register(rules.DiscountApprovalRule{})
	ruleEngine.Register(rules.CustomBuildApprovalRule{})

	if !ruleEngine.SetRuleEnabled("Custom Build Approval Rule", false) {
		t.Fatal("expected the CustomBuild Rule to be registered")
	}
	if ruleEngine.SetRuleEnabled("Unknown Rule", false) {
		t.Fatal("expected an unknown Rule name to be rejected")
	}

	if err := ruleEngine.ExecuteAll(memory); err != nil {
		t.Fatalf("execute rules: %v", err)
	}
	if len(memory.Findings) != 1 {
		t.Fatalf("expected one finding after disabling CustomBuild Rule, got %d", len(memory.Findings))
	}
	if memory.Findings[0].RuleName != "Discount Approval Rule" {
		t.Fatalf("expected discount finding, got %q", memory.Findings[0].RuleName)
	}
}
