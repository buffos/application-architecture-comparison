package rules

import (
	"fmt"

	"rules-engine-architecture/internal/engine"
)

// DiscountRejectionRule implements the PRD policy that discounts above 25%
// are always rejected.
type DiscountRejectionRule struct{}

func (DiscountRejectionRule) Name() string {
	return "Discount Rejection Rule"
}

func (DiscountRejectionRule) Priority() int {
	return 200
}

func (DiscountRejectionRule) ConflictGroup() string {
	return "discount-outcome"
}

func (DiscountRejectionRule) Evaluate(memory *engine.WorkingMemory) bool {
	return memory.Quote.DiscountPercent > 25
}

func (rule DiscountRejectionRule) Execute(memory *engine.WorkingMemory) error {
	memory.AddFinding(engine.Finding{
		RuleName: rule.Name(),
		Severity: "rejected",
		Message:  fmt.Sprintf("Discount of %d%% is always rejected", memory.Quote.DiscountPercent),
	})
	return nil
}
