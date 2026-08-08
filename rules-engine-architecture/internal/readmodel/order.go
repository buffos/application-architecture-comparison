package readmodel

import "rules-engine-architecture/internal/engine"

// OrderView is a read-only composition of lifecycle Facts and policy actions.
type OrderView struct {
	ID                    string
	CustomerID            string
	QuoteID               string
	Status                engine.OrderStatus
	PaymentStatus         engine.PaymentStatus
	ShipmentRequested     bool
	ShipmentAction        engine.ShipmentAction
	CancellationRequested bool
	CancellationAction    engine.CancellationAction
	ReturnAction          engine.ReturnAction
	Outcome               engine.DecisionOutcome
}

func ProjectOrder(
	memory *engine.WorkingMemory,
	decision engine.PolicyDecision,
) OrderView {
	return OrderView{
		ID:                    memory.Order.ID,
		CustomerID:            memory.Quote.CustomerID,
		QuoteID:               memory.Quote.ID,
		Status:                memory.Order.Status,
		PaymentStatus:         memory.Payment.Status,
		ShipmentRequested:     memory.Shipment.Requested,
		ShipmentAction:        decision.ShipmentAction,
		CancellationRequested: memory.Cancellation.Requested,
		CancellationAction:    decision.CancellationAction,
		ReturnAction:          decision.ReturnAction,
		Outcome:               decision.Outcome,
	}
}
