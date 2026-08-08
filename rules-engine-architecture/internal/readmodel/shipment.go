package readmodel

import "rules-engine-architecture/internal/engine"

// ShipmentView is the read contract consumed by a shipment workflow.
type ShipmentView struct {
	OrderID       string
	Requested     bool
	PaymentStatus engine.PaymentStatus
	InvoiceTerms  bool
	Action        engine.ShipmentAction
	Partial       bool
	Reason        string
}

func ProjectShipment(
	memory *engine.WorkingMemory,
	decision engine.PolicyDecision,
) ShipmentView {
	view := ShipmentView{
		OrderID:       memory.Order.ID,
		Requested:     memory.Shipment.Requested,
		PaymentStatus: memory.Payment.Status,
		InvoiceTerms:  memory.Customer.InvoiceTerms,
		Action:        decision.ShipmentAction,
		Partial:       decision.ShipmentAction == engine.ShipmentPartiallyAllowed,
	}

	for _, reason := range decision.Reasons {
		if reason.RuleName == "Shipment Payment Guard Rule" {
			view.Reason = reason.Message
			break
		}
	}

	return view
}
