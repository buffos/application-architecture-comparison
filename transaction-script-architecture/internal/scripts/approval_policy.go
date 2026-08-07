package scripts

import "transaction-script-architecture/internal/data"

const ApprovalReasonCustomBuild = "custom_build_requires_review"

// ApprovalDecision is a passive result value returned by a procedural rule
// evaluation. It intentionally does not know how a quote is persisted.
type ApprovalDecision struct {
	Required bool
	Reasons  []string
}

// EvaluateQuoteApproval evaluates the current quote data without changing it.
// Keeping this as a function, rather than an interface-backed policy object,
// keeps the Transaction Script boundary explicit.
func EvaluateQuoteApproval(quote data.Quote) ApprovalDecision {
	decision := ApprovalDecision{}
	for _, line := range quote.Lines {
		if line.ProductCategory != "CustomBuild" {
			continue
		}

		decision.Required = true
		decision.Reasons = append(decision.Reasons, ApprovalReasonCustomBuild)
		break
	}

	return decision
}
