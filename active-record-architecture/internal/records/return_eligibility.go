package records

const (
	ReturnEligibilityReasonClearance = "clearance_items_are_not_returnable"
)

// ReturnEligibilityDecision is the non-mutating result of return policy
// evaluation.
type ReturnEligibilityDecision struct {
	Eligible bool
	Reasons  []string
}

// EvaluateEligibility checks product-category policy against the snapshots
// carried by the return and its source order. It does not change either
// record.
func (request *ReturnRequest) EvaluateEligibility(order *Order) ReturnEligibilityDecision {
	decision := ReturnEligibilityDecision{Eligible: true}
	if request == nil || order == nil {
		return decision
	}

	for _, returnLine := range request.Lines {
		category := returnLine.ProductCategory
		if category == "" {
			for _, orderLine := range order.Lines {
				if orderLine.ID == returnLine.OrderLineID {
					category = orderLine.ProductCategory
					break
				}
			}
		}
		if category == "Clearance" {
			addReturnEligibilityReason(&decision, ReturnEligibilityReasonClearance)
		}
	}
	return decision
}

func addReturnEligibilityReason(decision *ReturnEligibilityDecision, reason string) {
	decision.Eligible = false
	for _, existing := range decision.Reasons {
		if existing == reason {
			return
		}
	}
	decision.Reasons = append(decision.Reasons, reason)
}
