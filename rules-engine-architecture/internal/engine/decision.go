package engine

type DecisionOutcome string

const (
	OutcomeAllowed       DecisionOutcome = "allowed"
	OutcomeNeedsApproval DecisionOutcome = "needs-approval"
	OutcomeRejected      DecisionOutcome = "rejected"
)

// PolicyDecision is the application-facing result of one Rule Engine cycle.
type PolicyDecision struct {
	Outcome DecisionOutcome
	Reasons []Finding
}

func (engine *Engine) Decide(memory *WorkingMemory) (PolicyDecision, error) {
	if err := engine.ExecuteAll(memory); err != nil {
		return PolicyDecision{}, err
	}

	return DecisionFromFindings(memory.Findings), nil
}

func DecisionFromFindings(findings []Finding) PolicyDecision {
	decision := PolicyDecision{
		Outcome: OutcomeAllowed,
		Reasons: append([]Finding(nil), findings...),
	}

	for _, finding := range findings {
		switch finding.Severity {
		case "rejected":
			decision.Outcome = OutcomeRejected
		case "approval-required":
			if decision.Outcome != OutcomeRejected {
				decision.Outcome = OutcomeNeedsApproval
			}
		}
	}

	return decision
}
