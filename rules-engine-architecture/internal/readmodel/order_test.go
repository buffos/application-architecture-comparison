package readmodel

import (
	"testing"

	"rules-engine-architecture/internal/engine"
)

func TestProjectOrderComposesLifecycleAndPolicyActions(t *testing.T) {
	memory := engine.NewWorkingMemory(
		engine.CustomerFact{ID: "CUST-001"},
		engine.QuoteFact{ID: "Q-1001", CustomerID: "CUST-001"},
		nil,
	)
	memory.Order = engine.OrderFact{ID: "ORD-1001", Status: engine.OrderConfirmed}
	memory.Payment = engine.PaymentFact{Status: engine.PaymentFailed}
	memory.Shipment = engine.ShipmentRequestFact{Requested: true}
	memory.Cancellation = engine.CancellationRequestFact{Requested: false}
	decision := engine.PolicyDecision{
		Outcome:            engine.OutcomeNeedsPaymentReview,
		ShipmentAction:     engine.ShipmentBlocked,
		CancellationAction: engine.CancellationNotRequested,
		ReturnAction:       engine.ReturnNotRequested,
	}

	view := ProjectOrder(memory, decision)

	if view.ID != "ORD-1001" || view.QuoteID != "Q-1001" {
		t.Fatalf("unexpected order identity: %+v", view)
	}
	if view.PaymentStatus != engine.PaymentFailed {
		t.Fatalf("expected failed payment in view, got %s", view.PaymentStatus)
	}
	if !view.ShipmentRequested || view.ShipmentAction != engine.ShipmentBlocked {
		t.Fatalf("expected requested blocked shipment, got %+v", view)
	}
	if view.Outcome != engine.OutcomeNeedsPaymentReview {
		t.Fatalf("expected payment review outcome, got %s", view.Outcome)
	}
}
