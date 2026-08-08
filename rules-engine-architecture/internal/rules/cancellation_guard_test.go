package rules

import (
	"testing"

	"rules-engine-architecture/internal/engine"
)

func TestCancellationGuardAllowsCancellationBeforeShipment(t *testing.T) {
	memory := cancellationMemory(engine.OrderConfirmed, true)
	rule := CancellationGuardRule{}

	if !rule.Evaluate(memory) {
		t.Fatal("expected a cancellation request to be evaluated")
	}
	if err := rule.Execute(memory); err != nil {
		t.Fatalf("execute rule: %v", err)
	}

	decision := engine.DecisionFromFindings(memory.Findings)
	if decision.CancellationAction != engine.CancellationAllowed {
		t.Fatalf("expected cancellation to be allowed, got %s", decision.CancellationAction)
	}
}

func TestCancellationGuardBlocksCancellationAfterShipment(t *testing.T) {
	memory := cancellationMemory(engine.OrderShipped, true)
	rule := CancellationGuardRule{}

	if err := rule.Execute(memory); err != nil {
		t.Fatalf("execute rule: %v", err)
	}

	decision := engine.DecisionFromFindings(memory.Findings)
	if decision.CancellationAction != engine.CancellationBlocked {
		t.Fatalf("expected cancellation to be blocked, got %s", decision.CancellationAction)
	}
}

func TestCancellationGuardIgnoresMissingRequest(t *testing.T) {
	memory := cancellationMemory(engine.OrderConfirmed, false)

	if (CancellationGuardRule{}).Evaluate(memory) {
		t.Fatal("expected no cancellation request to skip the Rule")
	}
}

func cancellationMemory(status engine.OrderStatus, requested bool) *engine.WorkingMemory {
	memory := engine.NewWorkingMemory(
		engine.CustomerFact{ID: "CUST-001"},
		engine.QuoteFact{ID: "Q-1001"},
		nil,
	)
	memory.Order = engine.OrderFact{ID: "ORD-1001", Status: status}
	memory.Cancellation = engine.CancellationRequestFact{Requested: requested}
	return memory
}
