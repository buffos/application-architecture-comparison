package scripts

import (
	"time"

	"transaction-script-architecture/internal/data"
)

const (
	ReturnEligibilityReasonClearance = "clearance_items_are_not_returnable"
	ReturnEligibilityReasonWindow    = "return_window_expired"
)

// ReturnEligibilityDecision is a passive result value for the procedural
// return policy.
type ReturnEligibilityDecision struct {
	Eligible bool
	Reasons  []string
}

// EvaluateReturnEligibility evaluates the current time-based policy.
func EvaluateReturnEligibility(order data.Order, request data.ReturnRequest) ReturnEligibilityDecision {
	return EvaluateReturnEligibilityAt(order, request, time.Now())
}

// EvaluateReturnEligibilityAt evaluates product-category and return-window
// rules without changing the order or return request.
func EvaluateReturnEligibilityAt(order data.Order, request data.ReturnRequest, now time.Time) ReturnEligibilityDecision {
	decision := ReturnEligibilityDecision{Eligible: true}
	for _, requestLine := range request.Lines {
		category := requestLine.ProductCategory
		windowDays := requestLine.ReturnWindowDays
		for _, orderLine := range order.Lines {
			if orderLine.ID != requestLine.OrderLineID {
				continue
			}
			if category == "" {
				category = orderLine.ProductCategory
			}
			if windowDays <= 0 {
				windowDays = orderLine.ReturnWindowDays
			}
			break
		}

		if category == "Clearance" {
			addReturnReason(&decision, ReturnEligibilityReasonClearance)
		}

		if !order.ShippedAt.IsZero() {
			if windowDays <= 0 {
				windowDays = data.DefaultReturnWindowDays
			}
			if now.After(order.ShippedAt.AddDate(0, 0, windowDays)) {
				addReturnReason(&decision, ReturnEligibilityReasonWindow)
			}
		}
	}

	return decision
}

func addReturnReason(decision *ReturnEligibilityDecision, reason string) {
	decision.Eligible = false
	for _, existing := range decision.Reasons {
		if existing == reason {
			return
		}
	}
	decision.Reasons = append(decision.Reasons, reason)
}
