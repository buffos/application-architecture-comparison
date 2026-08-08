package engine

type DecisionOutcome string

const (
	OutcomeAllowed            DecisionOutcome = "allowed"
	OutcomeNeedsApproval      DecisionOutcome = "needs-approval"
	OutcomeNeedsPaymentReview DecisionOutcome = "needs-payment-review"
	OutcomeNeedsReview        DecisionOutcome = "needs-review"
	OutcomeRejected           DecisionOutcome = "rejected"
)

type ReviewRequirement string

const (
	ReviewManagerApproval ReviewRequirement = "manager-approval"
	ReviewPayment         ReviewRequirement = "payment-review"
)

// PolicyDecision is the application-facing result of one Rule Engine cycle.
type PolicyDecision struct {
	Outcome         DecisionOutcome
	Reasons         []Finding
	RequiredReviews []ReviewRequirement
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
			decision.RequiredReviews = appendReviewRequirement(
				decision.RequiredReviews,
				ReviewManagerApproval,
			)
		case "payment-review":
			decision.RequiredReviews = appendReviewRequirement(
				decision.RequiredReviews,
				ReviewPayment,
			)
		}
	}

	if decision.Outcome == OutcomeRejected {
		decision.RequiredReviews = nil
	} else {
		switch len(decision.RequiredReviews) {
		case 1:
			switch decision.RequiredReviews[0] {
			case ReviewManagerApproval:
				decision.Outcome = OutcomeNeedsApproval
			case ReviewPayment:
				decision.Outcome = OutcomeNeedsPaymentReview
			}
		case 2:
			decision.Outcome = OutcomeNeedsReview
		}
	}

	return decision
}

func appendReviewRequirement(
	requirements []ReviewRequirement,
	requirement ReviewRequirement,
) []ReviewRequirement {
	for _, existing := range requirements {
		if existing == requirement {
			return requirements
		}
	}
	return append(requirements, requirement)
}
