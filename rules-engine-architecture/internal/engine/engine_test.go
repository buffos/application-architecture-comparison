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
	if len(memory.Trace) != 2 {
		t.Fatalf("expected two trace records, got %d", len(memory.Trace))
	}
	if memory.Trace[1].SkippedReason != "conflict group already resolved" {
		t.Fatalf("expected lower-priority Rule to be conflict-skipped, got %q", memory.Trace[1].SkippedReason)
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
				{ProductID: "PRD-002", Quantity: 1, UnitPriceCents: 125000},
			},
		},
		nil,
	)

	ruleEngine := engine.NewEngine()
	ruleEngine.Register(rules.DiscountApprovalRule{})
	ruleEngine.Register(rules.NewHighValuePaymentReviewRule(100000))

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

func TestEngineDecideReturnsApprovalWithReasons(t *testing.T) {
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

	decision, err := ruleEngine.Decide(memory)
	if err != nil {
		t.Fatalf("decide: %v", err)
	}
	if decision.Outcome != engine.OutcomeNeedsApproval {
		t.Fatalf("expected needs-approval, got %s", decision.Outcome)
	}
	if len(decision.Reasons) != 1 {
		t.Fatalf("expected the discount decision reason, got %d", len(decision.Reasons))
	}
	if len(decision.RequiredReviews) != 1 || decision.RequiredReviews[0] != engine.ReviewManagerApproval {
		t.Fatalf("expected manager approval requirement, got %v", decision.RequiredReviews)
	}
}

func TestEngineDecisionRejectsBeforeApproval(t *testing.T) {
	memory := engine.NewWorkingMemory(
		engine.CustomerFact{ID: "CUST-001"},
		engine.QuoteFact{
			ID:              "Q-1001",
			CustomerID:      "CUST-001",
			DiscountPercent: 30,
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
	ruleEngine.Register(rules.DiscountRejectionRule{})
	ruleEngine.Register(rules.CustomBuildApprovalRule{})

	decision, err := ruleEngine.Decide(memory)
	if err != nil {
		t.Fatalf("decide: %v", err)
	}
	if decision.Outcome != engine.OutcomeRejected {
		t.Fatalf("expected rejected, got %s", decision.Outcome)
	}
	if len(decision.Reasons) != 1 {
		t.Fatalf("expected only the rejection reason, got %d", len(decision.Reasons))
	}
	if len(decision.RequiredReviews) != 0 {
		t.Fatalf("expected no review requirements after rejection, got %v", decision.RequiredReviews)
	}
}

func TestEngineDecisionIncludesPaymentReviewRequirement(t *testing.T) {
	memory := engine.NewWorkingMemory(
		engine.CustomerFact{ID: "CUST-001"},
		engine.QuoteFact{
			ID:              "Q-1001",
			CustomerID:      "CUST-001",
			DiscountPercent: 20,
			Lines: []engine.QuoteLineFact{
				{ProductID: "PRD-002", Quantity: 1, UnitPriceCents: 125000},
			},
		},
		nil,
	)

	ruleEngine := engine.NewEngine()
	ruleEngine.Register(rules.DiscountApprovalRule{})
	ruleEngine.Register(rules.NewHighValuePaymentReviewRule(100000))

	decision, err := ruleEngine.Decide(memory)
	if err != nil {
		t.Fatalf("decide: %v", err)
	}
	if decision.Outcome != engine.OutcomeNeedsReview {
		t.Fatalf("expected needs-review, got %s", decision.Outcome)
	}
	if len(decision.RequiredReviews) != 2 {
		t.Fatalf("expected manager approval and payment review, got %v", decision.RequiredReviews)
	}
	if decision.RequiredReviews[1] != engine.ReviewPayment {
		t.Fatalf("expected payment review requirement, got %v", decision.RequiredReviews)
	}
}

func TestEngineRecordsRuleExecutionTrace(t *testing.T) {
	memory := engine.NewWorkingMemory(
		engine.CustomerFact{ID: "CUST-001"},
		engine.QuoteFact{
			ID:              "Q-1001",
			CustomerID:      "CUST-001",
			DiscountPercent: 20,
			Lines: []engine.QuoteLineFact{
				{ProductID: "PRD-002", Quantity: 1, UnitPriceCents: 125000},
			},
		},
		[]engine.ProductFact{
			{ID: "PRD-002", Category: "CustomBuild"},
		},
	)

	ruleEngine := engine.NewEngine()
	ruleEngine.Register(rules.DiscountRejectionRule{})
	ruleEngine.Register(rules.DiscountApprovalRule{})
	ruleEngine.Register(rules.CustomBuildApprovalRule{})
	ruleEngine.Register(rules.NewHighValuePaymentReviewRule(100000))

	if err := ruleEngine.ExecuteAll(memory); err != nil {
		t.Fatalf("execute rules: %v", err)
	}
	if len(memory.Trace) != 4 {
		t.Fatalf("expected four trace records, got %d", len(memory.Trace))
	}

	if memory.Trace[0].RuleName != "Discount Rejection Rule" || memory.Trace[0].Matched {
		t.Fatalf("expected first trace to be a non-matching rejection Rule, got %+v", memory.Trace[0])
	}
	if !memory.Trace[1].Executed || !memory.Trace[2].Executed || !memory.Trace[3].Executed {
		t.Fatalf("expected the remaining Rules to execute, got %+v", memory.Trace)
	}
}

func TestEngineChainsDerivedFactsUntilStable(t *testing.T) {
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

	ruleEngine := engine.NewEngine()
	// The gate has higher priority and therefore observes the derived Fact
	// only on the cycle after CustomBuildApprovalRule publishes it.
	ruleEngine.Register(rules.ApprovalWorkflowGateRule{})
	ruleEngine.Register(rules.CustomBuildApprovalRule{})

	decision, cycles, err := ruleEngine.DecideUntilStable(memory, 5)
	if err != nil {
		t.Fatalf("decide until stable: %v", err)
	}
	if cycles != 3 {
		t.Fatalf("expected two productive cycles plus a stable confirmation, got %d", cycles)
	}
	if !memory.HasDerivedFact(engine.ManagerApprovalRequiredFact) {
		t.Fatal("expected the CustomBuild Rule to publish a derived approval Fact")
	}
	if len(memory.Findings) != 1 {
		t.Fatalf("expected only the consumer workflow finding, got %d", len(memory.Findings))
	}
	if memory.Findings[0].RuleName != "Approval Workflow Gate Rule" {
		t.Fatalf("expected the workflow gate to add the finding, got %q", memory.Findings[0].RuleName)
	}
	if decision.Outcome != engine.OutcomeNeedsApproval {
		t.Fatalf("expected needs-approval, got %s", decision.Outcome)
	}
}

func TestEngineRecomputesAfterSourceFactChange(t *testing.T) {
	memory := engine.NewWorkingMemory(
		engine.CustomerFact{ID: "CUST-001"},
		engine.QuoteFact{
			ID: "Q-1001",
			Lines: []engine.QuoteLineFact{
				{ProductID: "PRD-002", Quantity: 1},
			},
		},
		[]engine.ProductFact{
			{ID: "PRD-001", Category: "Standard"},
			{ID: "PRD-002", Category: "CustomBuild"},
		},
	)

	ruleEngine := engine.NewEngine()
	ruleEngine.Register(rules.ApprovalWorkflowGateRule{})
	ruleEngine.Register(rules.CustomBuildApprovalRule{})

	decision, cycles, err := ruleEngine.RecomputeDecision(memory, 5)
	if err != nil {
		t.Fatalf("initial recompute: %v", err)
	}
	if cycles != 3 {
		t.Fatalf("expected three cycles for the initial inference, got %d", cycles)
	}
	if decision.Outcome != engine.OutcomeNeedsApproval {
		t.Fatalf("expected initial needs-approval decision, got %s", decision.Outcome)
	}
	if !memory.HasDerivedFact(engine.ManagerApprovalRequiredFact) {
		t.Fatal("expected the initial CustomBuild inference")
	}

	memory.Quote.Lines = []engine.QuoteLineFact{
		{ProductID: "PRD-001", Quantity: 1},
	}

	decision, cycles, err = ruleEngine.RecomputeDecision(memory, 5)
	if err != nil {
		t.Fatalf("recompute after quote change: %v", err)
	}
	if cycles != 1 {
		t.Fatalf("expected one stable cycle after removing CustomBuild, got %d", cycles)
	}
	if decision.Outcome != engine.OutcomeAllowed {
		t.Fatalf("expected allowed decision after recomputation, got %s", decision.Outcome)
	}
	if memory.HasDerivedFact(engine.ManagerApprovalRequiredFact) {
		t.Fatal("expected the stale approval fact to be removed by recomputation")
	}
	if len(memory.Findings) != 0 {
		t.Fatalf("expected no stale findings after recomputation, got %d", len(memory.Findings))
	}
	if len(memory.Trace) != 2 {
		t.Fatalf("expected the trace to contain only the fresh cycle, got %d records", len(memory.Trace))
	}
}
