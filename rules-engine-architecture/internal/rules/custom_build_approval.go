package rules

import (
	"rules-engine-architecture/internal/engine"
)

// CustomBuildApprovalRule implements the PRD policy that CustomBuild products
// require approval before order conversion.
type CustomBuildApprovalRule struct{}

func (CustomBuildApprovalRule) Name() string {
	return "Custom Build Approval Rule"
}

func (CustomBuildApprovalRule) Priority() int {
	return 100
}

func (CustomBuildApprovalRule) ConflictGroup() string {
	// This finding complements discount outcomes instead of competing with them.
	return ""
}

func (CustomBuildApprovalRule) Evaluate(memory *engine.WorkingMemory) bool {
	for _, line := range memory.Quote.Lines {
		for _, product := range memory.Products {
			if line.ProductID == product.ID && product.Category == "CustomBuild" {
				return true
			}
		}
	}
	return false
}

func (CustomBuildApprovalRule) Execute(memory *engine.WorkingMemory) error {
	// This Rule is a fact producer. The workflow consequence is intentionally
	// left to ApprovalWorkflowGateRule, which consumes the derived Fact.
	memory.AddDerivedFact(engine.DerivedFact{
		Name:  engine.ManagerApprovalRequiredFact,
		Value: memory.Quote.ID,
	})
	return nil
}
