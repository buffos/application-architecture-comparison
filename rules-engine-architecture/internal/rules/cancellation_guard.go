package rules

import "rules-engine-architecture/internal/engine"

// CancellationGuardRule allows cancellation before shipment and directs
// shipped orders to the return flow.
type CancellationGuardRule struct{}

func (CancellationGuardRule) Name() string {
	return "Cancellation Guard Rule"
}

func (CancellationGuardRule) Priority() int {
	return 60
}

func (CancellationGuardRule) ConflictGroup() string {
	return ""
}

func (CancellationGuardRule) Evaluate(memory *engine.WorkingMemory) bool {
	return memory.Cancellation.Requested
}

func (rule CancellationGuardRule) Execute(memory *engine.WorkingMemory) error {
	if memory.Order.Status == engine.OrderShipped {
		memory.AddFinding(engine.Finding{
			RuleName: rule.Name(),
			Severity: "cancellation-blocked",
			Message:  "Cancellation is blocked after shipment; use the return flow",
		})
		return nil
	}

	memory.AddFinding(engine.Finding{
		RuleName: rule.Name(),
		Severity: "cancellation-allowed",
		Message:  "Cancellation is allowed before shipment",
	})
	return nil
}
