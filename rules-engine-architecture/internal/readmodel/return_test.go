package readmodel

import (
	"testing"

	"rules-engine-architecture/internal/engine"
	"rules-engine-architecture/internal/rules"
)

func TestProjectReturnBuildsIndependentReadView(t *testing.T) {
	memory := engine.NewWorkingMemory(
		engine.CustomerFact{ID: "CUST-001"},
		engine.QuoteFact{ID: "Q-1001"},
		[]engine.ProductFact{{ID: "PRD-002", Category: "Standard"}},
	)
	memory.Order = engine.OrderFact{ID: "ORD-1001", Status: engine.OrderShipped}
	memory.ReturnRequest = engine.ReturnRequestFact{
		Requested:         true,
		ProductID:         "PRD-002",
		Quantity:          1,
		ShippedQuantity:   2,
		DaysSinceShipment: 5,
		ReturnWindowDays:  30,
		RequestedBy: engine.ActorFact{
			ID:   "warehouse-clerk-001",
			Role: "warehouse-clerk",
		},
	}

	rule := rules.ReturnPolicyRule{}
	if err := rule.Execute(memory); err != nil {
		t.Fatalf("execute Rule: %v", err)
	}
	decision := engine.DecisionFromFindings(memory.Findings)
	view := ProjectReturn(memory, decision)

	if view.Action != engine.ReturnAllowed {
		t.Fatalf("expected allowed return view, got %s", view.Action)
	}
	if view.RemainingQuantity != 2 {
		t.Fatalf("expected two remaining units, got %d", view.RemainingQuantity)
	}
	if view.RequesterID != "warehouse-clerk-001" {
		t.Fatalf("expected requester in view, got %q", view.RequesterID)
	}
	if view.Reason == "" {
		t.Fatal("expected the Rule explanation in the view")
	}

	memory.ReturnRequest.Quantity = 99
	view.RequestedQuantity = 7
	if view.RequestedQuantity != 7 {
		t.Fatal("expected the view to be independently mutable")
	}
	if memory.ReturnRequest.Quantity != 99 {
		t.Fatal("expected the source Fact to remain independently mutable")
	}
}
