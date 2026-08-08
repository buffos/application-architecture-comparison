package engine_test

import (
	"testing"

	"rules-engine-architecture/internal/engine"
	"rules-engine-architecture/internal/rules"
)

func TestPaymentReviewApprovalClearsFindingAfterRecompute(t *testing.T) {
	memory := paymentReviewMemory()
	ruleEngine := engine.NewEngine()
	ruleEngine.Register(rules.NewHighValuePaymentReviewRule(100000))

	decision, _, err := ruleEngine.DecideUntilStable(memory, 5)
	if err != nil {
		t.Fatalf("initial decision: %v", err)
	}
	if decision.Outcome != engine.OutcomeNeedsPaymentReview {
		t.Fatalf("expected payment review, got %s", decision.Outcome)
	}

	memory.PaymentReview = engine.PaymentReviewFact{
		Status:     engine.PaymentReviewApproved,
		ReviewedBy: "payment-manager",
	}
	decision, cycles, err := ruleEngine.RecomputeDecision(memory, 5)
	if err != nil {
		t.Fatalf("decision after approval: %v", err)
	}
	if cycles != 1 {
		t.Fatalf("expected one stable cycle after approval, got %d", cycles)
	}
	if decision.Outcome != engine.OutcomeAllowed || len(decision.RequiredReviews) != 0 {
		t.Fatalf("expected payment review to clear, got %+v", decision)
	}
}

func TestPaymentReviewRejectionIsTerminal(t *testing.T) {
	memory := paymentReviewMemory()
	memory.PaymentReview.Status = engine.PaymentReviewRejected
	ruleEngine := engine.NewEngine()
	ruleEngine.Register(rules.NewHighValuePaymentReviewRule(100000))

	decision, err := ruleEngine.Decide(memory)
	if err != nil {
		t.Fatalf("rejected review decision: %v", err)
	}
	if decision.Outcome != engine.OutcomeRejected {
		t.Fatalf("expected rejected outcome, got %s", decision.Outcome)
	}
}

func paymentReviewMemory() *engine.WorkingMemory {
	return engine.NewWorkingMemory(
		engine.CustomerFact{ID: "CUST-001"},
		engine.QuoteFact{
			ID: "Q-1001",
			Lines: []engine.QuoteLineFact{
				{ProductID: "PRD-002", Quantity: 1, UnitPriceCents: 125000},
			},
		},
		nil,
	)
}
