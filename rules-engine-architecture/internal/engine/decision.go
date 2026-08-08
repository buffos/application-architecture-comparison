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

type FulfillmentAction string

const (
	FulfillmentNone      FulfillmentAction = "none"
	FulfillmentBackorder FulfillmentAction = "backorder"
	FulfillmentReject    FulfillmentAction = "reject"
)

type ShipmentAction string

const (
	ShipmentNotRequested ShipmentAction = "not-requested"
	ShipmentAllowed      ShipmentAction = "allowed"
	ShipmentBlocked      ShipmentAction = "blocked"
)

type CancellationAction string

const (
	CancellationNotRequested CancellationAction = "not-requested"
	CancellationAllowed      CancellationAction = "allowed"
	CancellationBlocked      CancellationAction = "blocked"
)

// PolicyDecision is the application-facing result of one Rule Engine cycle.
type PolicyDecision struct {
	Outcome            DecisionOutcome
	FulfillmentAction  FulfillmentAction
	ShipmentAction     ShipmentAction
	CancellationAction CancellationAction
	Reasons            []Finding
	RequiredReviews    []ReviewRequirement
}

func (engine *Engine) Decide(memory *WorkingMemory) (PolicyDecision, error) {
	if err := engine.ExecuteAll(memory); err != nil {
		return PolicyDecision{}, err
	}

	return DecisionFromFindings(memory.Findings), nil
}

func (engine *Engine) DecideUntilStable(
	memory *WorkingMemory,
	maxCycles int,
) (PolicyDecision, int, error) {
	cycles, err := engine.ExecuteUntilStable(memory, maxCycles)
	if err != nil {
		return PolicyDecision{}, cycles, err
	}

	return DecisionFromFindings(memory.Findings), cycles, nil
}

// RecomputeDecision starts a fresh inference session for the current source
// Facts and returns the resulting policy decision.
func (engine *Engine) RecomputeDecision(
	memory *WorkingMemory,
	maxCycles int,
) (PolicyDecision, int, error) {
	cycles, err := engine.RecomputeUntilStable(memory, maxCycles)
	if err != nil {
		return PolicyDecision{}, cycles, err
	}

	return DecisionFromFindings(memory.Findings), cycles, nil
}

func DecisionFromFindings(findings []Finding) PolicyDecision {
	decision := PolicyDecision{
		Outcome:            OutcomeAllowed,
		FulfillmentAction:  FulfillmentNone,
		ShipmentAction:     ShipmentNotRequested,
		CancellationAction: CancellationNotRequested,
		Reasons:            append([]Finding(nil), findings...),
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
		case "conversion-blocked":
			decision.RequiredReviews = appendReviewRequirement(
				decision.RequiredReviews,
				ReviewManagerApproval,
			)
		case "inventory-backorder":
			if decision.FulfillmentAction != FulfillmentReject {
				decision.FulfillmentAction = FulfillmentBackorder
			}
		case "inventory-rejected":
			decision.FulfillmentAction = FulfillmentReject
			decision.Outcome = OutcomeRejected
		case "shipment-allowed":
			decision.ShipmentAction = ShipmentAllowed
		case "shipment-blocked":
			decision.ShipmentAction = ShipmentBlocked
		case "cancellation-allowed":
			decision.CancellationAction = CancellationAllowed
		case "cancellation-blocked":
			decision.CancellationAction = CancellationBlocked
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
