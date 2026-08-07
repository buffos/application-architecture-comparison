package rules

import (
	"fmt"

	"rules-engine-architecture/internal/engine"
)

// DiscountApprovalRule implements the PRD policy that discounts above 15%
// require manager approval.
type DiscountApprovalRule struct{}

func (DiscountApprovalRule) Name() string {
	return "Discount Approval Rule"
}

func (DiscountApprovalRule) Priority() int {
	return 100
}

func (DiscountApprovalRule) Evaluate(memory *engine.WorkingMemory) bool {
	return memory.Quote.DiscountPercent > 15
}

func (rule DiscountApprovalRule) Execute(memory *engine.WorkingMemory) error {
	memory.AddFinding(engine.Finding{
		RuleName: rule.Name(),
		Severity: "approval-required",
		Message:  fmt.Sprintf("Discount of %d%% requires manager approval", memory.Quote.DiscountPercent),
	})
	return nil
}
