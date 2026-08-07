package scripts

import "transaction-script-architecture/internal/data"

const ReturnEligibilityReasonClearance = "clearance_items_are_not_returnable"

// ReturnEligibilityDecision is a passive result value for the procedural
// return policy.
type ReturnEligibilityDecision struct {
	Eligible bool
	Reasons  []string
}

// EvaluateReturnEligibility evaluates product-category rules without changing
// the order or return request.
func EvaluateReturnEligibility(order data.Order, request data.ReturnRequest) ReturnEligibilityDecision {
	decision := ReturnEligibilityDecision{Eligible: true}
	for _, requestLine := range request.Lines {
		category := requestLine.ProductCategory
		if category == "" {
			for _, orderLine := range order.Lines {
				if orderLine.ID == requestLine.OrderLineID {
					category = orderLine.ProductCategory
					break
				}
			}
		}

		if category == "Clearance" {
			decision.Eligible = false
			decision.Reasons = append(decision.Reasons, ReturnEligibilityReasonClearance)
		}
	}

	return decision
}
