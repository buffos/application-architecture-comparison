package rules

import "rules-engine-architecture/internal/engine"

// ShipmentPaymentGuardRule decides whether a requested shipment is allowed by
// payment state or by the customer's invoice terms.
type ShipmentPaymentGuardRule struct{}

func (ShipmentPaymentGuardRule) Name() string {
	return "Shipment Payment Guard Rule"
}

func (ShipmentPaymentGuardRule) Priority() int {
	return 70
}

func (ShipmentPaymentGuardRule) ConflictGroup() string {
	return ""
}

func (ShipmentPaymentGuardRule) Evaluate(memory *engine.WorkingMemory) bool {
	return memory.Shipment.Requested
}

func (rule ShipmentPaymentGuardRule) Execute(memory *engine.WorkingMemory) error {
	if memory.Payment.Status == engine.PaymentAccepted || memory.Customer.InvoiceTerms {
		memory.AddFinding(engine.Finding{
			RuleName: rule.Name(),
			Severity: "shipment-allowed",
			Message:  "Shipment is allowed because payment is accepted or invoice terms apply",
		})
		return nil
	}

	memory.AddFinding(engine.Finding{
		RuleName: rule.Name(),
		Severity: "shipment-blocked",
		Message:  "Shipment is blocked until payment is accepted or invoice terms apply",
	})
	return nil
}
