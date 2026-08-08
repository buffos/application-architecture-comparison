package rules

import "rules-engine-architecture/internal/engine"

// PartialShipmentRule identifies a permitted shipment where at least one line
// has shipped some, but not all, of its ordered quantity.
type PartialShipmentRule struct{}

func (PartialShipmentRule) Name() string {
	return "Partial Shipment Rule"
}

func (PartialShipmentRule) Priority() int {
	return 65
}

func (PartialShipmentRule) ConflictGroup() string {
	return ""
}

func (PartialShipmentRule) Evaluate(memory *engine.WorkingMemory) bool {
	if !memory.Shipment.Requested {
		return false
	}
	if memory.Payment.Status != engine.PaymentAccepted && !memory.Customer.InvoiceTerms {
		return false
	}

	for _, line := range memory.Shipment.Lines {
		if line.AlreadyShippedQuantity > 0 &&
			line.AlreadyShippedQuantity < line.OrderedQuantity {
			return true
		}
	}
	return false
}

func (rule PartialShipmentRule) Execute(memory *engine.WorkingMemory) error {
	memory.AddFinding(engine.Finding{
		RuleName: rule.Name(),
		Severity: "shipment-partial",
		Message:  "Shipment is allowed for the remaining quantities after a partial shipment",
	})
	return nil
}
