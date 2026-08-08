package rules

import (
	"testing"

	"rules-engine-architecture/internal/engine"
)

func TestReturnPolicyAllowsEligibleReturn(t *testing.T) {
	memory := returnMemory(engine.OrderShipped, "Standard", engine.ReturnRequestFact{
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
	})
	rule := ReturnPolicyRule{}

	if !rule.Evaluate(memory) {
		t.Fatal("expected a return request to be evaluated")
	}
	if err := rule.Execute(memory); err != nil {
		t.Fatalf("execute rule: %v", err)
	}

	decision := engine.DecisionFromFindings(memory.Findings)
	if decision.ReturnAction != engine.ReturnAllowed {
		t.Fatalf("expected return to be allowed, got %s", decision.ReturnAction)
	}
	if decision.Outcome != engine.OutcomeAllowed {
		t.Fatalf("expected allowed outcome, got %s", decision.Outcome)
	}
}

func TestReturnPolicyRejectsUnshippedOrder(t *testing.T) {
	memory := returnMemory(engine.OrderConfirmed, "Standard", validReturnRequest())

	if err := (ReturnPolicyRule{}).Execute(memory); err != nil {
		t.Fatalf("execute rule: %v", err)
	}

	decision := engine.DecisionFromFindings(memory.Findings)
	if decision.ReturnAction != engine.ReturnRejected {
		t.Fatalf("expected return to be rejected, got %s", decision.ReturnAction)
	}
	if decision.Outcome != engine.OutcomeRejected {
		t.Fatalf("expected rejected outcome, got %s", decision.Outcome)
	}
}

func TestReturnPolicyAllowsPartialReturnByLine(t *testing.T) {
	memory := engine.NewWorkingMemory(
		engine.CustomerFact{ID: "CUST-001"},
		engine.QuoteFact{ID: "Q-1001"},
		[]engine.ProductFact{
			{ID: "PRD-001", Category: "Standard"},
			{ID: "PRD-002", Category: "Standard"},
		},
	)
	memory.Order = engine.OrderFact{ID: "ORD-1001", Status: engine.OrderShipped}
	memory.ReturnRequest = engine.ReturnRequestFact{
		Requested:         true,
		DaysSinceShipment: 5,
		ReturnWindowDays:  30,
		RequestedBy: engine.ActorFact{
			ID:   "warehouse-clerk-001",
			Role: "warehouse-clerk",
		},
		Lines: []engine.ReturnLineFact{
			{ProductID: "PRD-001", Quantity: 1, ShippedQuantity: 1},
			{ProductID: "PRD-002", Quantity: 2, ShippedQuantity: 1},
		},
	}

	if err := (ReturnPolicyRule{}).Execute(memory); err != nil {
		t.Fatalf("execute rule: %v", err)
	}

	decision := engine.DecisionFromFindings(memory.Findings)
	if decision.ReturnAction != engine.ReturnPartiallyAllowed {
		t.Fatalf("expected partial return action, got %s", decision.ReturnAction)
	}
	if decision.Outcome != engine.OutcomeAllowed {
		t.Fatalf("expected partial return to remain allowed, got %s", decision.Outcome)
	}
}

func TestReturnPolicyRejectsMissingActor(t *testing.T) {
	request := validReturnRequest()
	request.RequestedBy = engine.ActorFact{}
	memory := returnMemory(engine.OrderShipped, "Standard", request)

	if err := (ReturnPolicyRule{}).Execute(memory); err != nil {
		t.Fatalf("execute rule: %v", err)
	}

	if memory.Findings[0].Severity != "return-rejected" {
		t.Fatalf("expected return rejection finding, got %+v", memory.Findings)
	}
}

func TestReturnPolicyRejectsClearanceProduct(t *testing.T) {
	memory := returnMemory(engine.OrderShipped, "Clearance", validReturnRequest())

	if err := (ReturnPolicyRule{}).Execute(memory); err != nil {
		t.Fatalf("execute rule: %v", err)
	}

	if memory.Findings[0].Severity != "return-rejected" {
		t.Fatalf("expected return rejection finding, got %+v", memory.Findings)
	}
}

func TestReturnPolicyRejectsLateReturn(t *testing.T) {
	request := validReturnRequest()
	request.DaysSinceShipment = 31
	memory := returnMemory(engine.OrderShipped, "Standard", request)

	if err := (ReturnPolicyRule{}).Execute(memory); err != nil {
		t.Fatalf("execute rule: %v", err)
	}

	if memory.Findings[0].Severity != "return-rejected" {
		t.Fatalf("expected return rejection finding, got %+v", memory.Findings)
	}
}

func TestReturnPolicyRejectsQuantityAboveRemainingAmount(t *testing.T) {
	request := validReturnRequest()
	request.Quantity = 2
	request.PreviouslyReturnedQuantity = 1
	memory := returnMemory(engine.OrderShipped, "Standard", request)

	if err := (ReturnPolicyRule{}).Execute(memory); err != nil {
		t.Fatalf("execute rule: %v", err)
	}

	if memory.Findings[0].Severity != "return-rejected" {
		t.Fatalf("expected return rejection finding, got %+v", memory.Findings)
	}
}

func TestReturnPolicyIgnoresMissingRequest(t *testing.T) {
	memory := returnMemory(engine.OrderShipped, "Standard", engine.ReturnRequestFact{})

	if (ReturnPolicyRule{}).Evaluate(memory) {
		t.Fatal("expected missing return request to skip the Rule")
	}
}

func validReturnRequest() engine.ReturnRequestFact {
	return engine.ReturnRequestFact{
		Requested:         true,
		ProductID:         "PRD-002",
		Quantity:          1,
		ShippedQuantity:   2,
		ReturnWindowDays:  30,
		DaysSinceShipment: 5,
		RequestedBy: engine.ActorFact{
			ID:   "warehouse-clerk-001",
			Role: "warehouse-clerk",
		},
	}
}

func returnMemory(
	status engine.OrderStatus,
	category string,
	request engine.ReturnRequestFact,
) *engine.WorkingMemory {
	memory := engine.NewWorkingMemory(
		engine.CustomerFact{ID: "CUST-001"},
		engine.QuoteFact{ID: "Q-1001"},
		[]engine.ProductFact{{ID: "PRD-002", Category: category}},
	)
	memory.Order = engine.OrderFact{ID: "ORD-1001", Status: status}
	memory.ReturnRequest = request
	return memory
}
