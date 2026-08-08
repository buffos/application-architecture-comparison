package rules

import (
	"testing"

	"rules-engine-architecture/internal/engine"
)

func TestShipmentPaymentGuardAllowsAcceptedPayment(t *testing.T) {
	memory := shipmentMemory(engine.PaymentAccepted, false, true)
	rule := ShipmentPaymentGuardRule{}

	if !rule.Evaluate(memory) {
		t.Fatal("expected a requested shipment to be evaluated")
	}
	if err := rule.Execute(memory); err != nil {
		t.Fatalf("execute rule: %v", err)
	}

	decision := engine.DecisionFromFindings(memory.Findings)
	if decision.ShipmentAction != engine.ShipmentAllowed {
		t.Fatalf("expected allowed shipment, got %s", decision.ShipmentAction)
	}
}

func TestShipmentPaymentGuardBlocksUnacceptedPayment(t *testing.T) {
	memory := shipmentMemory(engine.PaymentFailed, false, true)
	rule := ShipmentPaymentGuardRule{}

	if err := rule.Execute(memory); err != nil {
		t.Fatalf("execute rule: %v", err)
	}

	decision := engine.DecisionFromFindings(memory.Findings)
	if decision.ShipmentAction != engine.ShipmentBlocked {
		t.Fatalf("expected blocked shipment, got %s", decision.ShipmentAction)
	}
}

func TestShipmentPaymentGuardAllowsInvoiceTerms(t *testing.T) {
	memory := shipmentMemory(engine.PaymentFailed, true, true)
	rule := ShipmentPaymentGuardRule{}

	if err := rule.Execute(memory); err != nil {
		t.Fatalf("execute rule: %v", err)
	}

	decision := engine.DecisionFromFindings(memory.Findings)
	if decision.ShipmentAction != engine.ShipmentAllowed {
		t.Fatalf("expected invoice terms to allow shipment, got %s", decision.ShipmentAction)
	}
}

func TestShipmentPaymentGuardIgnoresUnrequestedShipment(t *testing.T) {
	memory := shipmentMemory(engine.PaymentFailed, false, false)
	rule := ShipmentPaymentGuardRule{}

	if rule.Evaluate(memory) {
		t.Fatal("expected an unrequested shipment not to match")
	}
}

func shipmentMemory(
	status engine.PaymentStatus,
	invoiceTerms bool,
	requested bool,
) *engine.WorkingMemory {
	memory := engine.NewWorkingMemory(
		engine.CustomerFact{ID: "CUST-001", InvoiceTerms: invoiceTerms},
		engine.QuoteFact{ID: "Q-1001"},
		nil,
	)
	memory.Payment = engine.PaymentFact{Status: status}
	memory.Shipment = engine.ShipmentRequestFact{Requested: requested}
	return memory
}
