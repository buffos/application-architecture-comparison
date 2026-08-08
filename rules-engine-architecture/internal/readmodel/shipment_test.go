package readmodel

import (
	"testing"

	"rules-engine-architecture/internal/engine"
)

func TestProjectShipmentBuildsFocusedBlockedView(t *testing.T) {
	memory := engine.NewWorkingMemory(
		engine.CustomerFact{ID: "CUST-001", InvoiceTerms: false},
		engine.QuoteFact{ID: "Q-1001"},
		nil,
	)
	memory.Order = engine.OrderFact{ID: "ORD-1001", Status: engine.OrderConfirmed}
	memory.Payment = engine.PaymentFact{Status: engine.PaymentFailed}
	memory.Shipment = engine.ShipmentRequestFact{Requested: true}
	decision := engine.PolicyDecision{
		ShipmentAction: engine.ShipmentBlocked,
		Reasons: []engine.Finding{{
			RuleName: "Shipment Payment Guard Rule",
			Message:  "payment is not accepted",
		}},
	}

	view := ProjectShipment(memory, decision)

	if view.OrderID != "ORD-1001" || !view.Requested {
		t.Fatalf("unexpected shipment identity/request: %+v", view)
	}
	if view.Action != engine.ShipmentBlocked {
		t.Fatalf("expected blocked shipment, got %s", view.Action)
	}
	if view.Reason != "payment is not accepted" {
		t.Fatalf("expected payment reason, got %q", view.Reason)
	}
}
