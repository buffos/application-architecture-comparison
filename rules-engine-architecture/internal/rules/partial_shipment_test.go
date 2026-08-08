package rules

import (
	"testing"

	"rules-engine-architecture/internal/engine"
)

func TestPartialShipmentRuleMarksAllowedPartialShipment(t *testing.T) {
	memory := shipmentMemory(engine.PaymentAccepted, false, true)
	memory.Shipment.Lines = []engine.ShipmentLineFact{{
		ProductID:              "PRD-002",
		OrderedQuantity:        3,
		AlreadyShippedQuantity: 1,
	}}
	rule := PartialShipmentRule{}

	if !rule.Evaluate(memory) {
		t.Fatal("expected a permitted partial shipment to match")
	}
	if err := rule.Execute(memory); err != nil {
		t.Fatalf("execute Rule: %v", err)
	}

	decision := engine.DecisionFromFindings(memory.Findings)
	if decision.ShipmentAction != engine.ShipmentPartiallyAllowed {
		t.Fatalf("expected partial shipment action, got %s", decision.ShipmentAction)
	}
}

func TestPartialShipmentRuleDoesNotOverridePaymentBlock(t *testing.T) {
	memory := shipmentMemory(engine.PaymentFailed, false, true)
	memory.Shipment.Lines = []engine.ShipmentLineFact{{
		ProductID:              "PRD-002",
		OrderedQuantity:        3,
		AlreadyShippedQuantity: 1,
	}}

	if (PartialShipmentRule{}).Evaluate(memory) {
		t.Fatal("expected payment block to prevent partial shipment action")
	}
}

func TestPartialShipmentRuleIgnoresFullyUnshippedLine(t *testing.T) {
	memory := shipmentMemory(engine.PaymentAccepted, false, true)
	memory.Shipment.Lines = []engine.ShipmentLineFact{{
		ProductID:              "PRD-002",
		OrderedQuantity:        3,
		AlreadyShippedQuantity: 0,
	}}

	if (PartialShipmentRule{}).Evaluate(memory) {
		t.Fatal("expected a full first shipment not to be classified as partial")
	}
}
