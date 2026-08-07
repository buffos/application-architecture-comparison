package records

import "time"

const (
	ReturnEligibilityReasonClearance = "clearance_items_are_not_returnable"
	ReturnEligibilityReasonWindow    = "return_window_expired"
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
	return request.EvaluateEligibilityAt(order, time.Now())
}

// EvaluateEligibilityAt evaluates category and date-window policy without
// changing either record.
func (request *ReturnRequest) EvaluateEligibilityAt(order *Order, now time.Time) ReturnEligibilityDecision {
	decision := ReturnEligibilityDecision{Eligible: true}
	if request == nil || order == nil {
		return decision
	}

	for _, returnLine := range request.Lines {
		category := returnLine.ProductCategory
		windowDays := returnLine.ReturnWindowDays
		if category == "" {
			for _, orderLine := range order.Lines {
				if orderLine.ID == returnLine.OrderLineID {
					category = orderLine.ProductCategory
					if windowDays <= 0 {
						windowDays = orderLine.ReturnWindowDays
					}
					break
				}
			}
		}
		if category == "Clearance" {
			addReturnEligibilityReason(&decision, ReturnEligibilityReasonClearance)
		}
		if !order.ShippedAt.IsZero() {
			if windowDays <= 0 {
				windowDays = DefaultReturnWindowDays
			}
			if now.After(order.ShippedAt.AddDate(0, 0, windowDays)) {
				addReturnEligibilityReason(&decision, ReturnEligibilityReasonWindow)
			}
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
