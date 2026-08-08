package rules

import "rules-engine-architecture/internal/engine"

// ApprovalWorkflowGateRule consumes a derived approval fact and exposes the
// workflow consequence to the policy decision.
type ApprovalWorkflowGateRule struct{}

func (ApprovalWorkflowGateRule) Name() string {
	return "Approval Workflow Gate Rule"
}

func (ApprovalWorkflowGateRule) Priority() int {
	return 150
}

func (ApprovalWorkflowGateRule) ConflictGroup() string {
	return ""
}

func (ApprovalWorkflowGateRule) Evaluate(memory *engine.WorkingMemory) bool {
	return memory.ManagerApproval.Status != engine.ApprovalApproved &&
		memory.HasDerivedFact(engine.ManagerApprovalRequiredFact)
}

func (rule ApprovalWorkflowGateRule) Execute(memory *engine.WorkingMemory) error {
	memory.AddFinding(engine.Finding{
		RuleName: rule.Name(),
		Severity: "conversion-blocked",
		Message:  "Quote conversion is blocked until manager approval",
	})
	return nil
}
