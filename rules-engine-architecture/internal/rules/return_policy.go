package rules

import (
	"fmt"
	"strings"

	"rules-engine-architecture/internal/engine"
)

// ReturnPolicyRule evaluates eligibility without applying inventory or refund
// side effects.
type ReturnPolicyRule struct{}

func (ReturnPolicyRule) Name() string {
	return "Return Policy Rule"
}

func (ReturnPolicyRule) Priority() int {
	return 65
}

func (ReturnPolicyRule) ConflictGroup() string {
	return ""
}

func (ReturnPolicyRule) Evaluate(memory *engine.WorkingMemory) bool {
	return memory.ReturnRequest.Requested
}

func (rule ReturnPolicyRule) Execute(memory *engine.WorkingMemory) error {
	request := memory.ReturnRequest
	if request.RequestedBy.ID == "" || request.RequestedBy.Role == "" {
		returnReject(memory, "Return request must include the requester identity and role")
		return nil
	}

	if memory.Order.Status != engine.OrderShipped {
		returnReject(memory, "A return requires an order that has already shipped")
		return nil
	}

	if request.DaysSinceShipment > request.ReturnWindowDays {
		returnReject(memory, fmt.Sprintf(
			"Return is %d days after shipment, beyond the %d-day return window",
			request.DaysSinceShipment,
			request.ReturnWindowDays,
		))
		return nil
	}

	acceptedProducts := make([]string, 0)
	rejectedReasons := make([]string, 0)
	for _, line := range returnLines(request) {
		product, found := productForReturn(memory, line.ProductID)
		if !found {
			rejectedReasons = append(rejectedReasons, fmt.Sprintf(
				"unknown product %s",
				line.ProductID,
			))
			continue
		}

		if product.Category == "Clearance" {
			rejectedReasons = append(rejectedReasons, fmt.Sprintf(
				"product %s is clearance",
				product.ID,
			))
			continue
		}

		if line.Quantity <= 0 {
			rejectedReasons = append(rejectedReasons, fmt.Sprintf(
				"product %s has a non-positive return quantity",
				product.ID,
			))
			continue
		}

		remainingQuantity := line.ShippedQuantity - line.PreviouslyReturnedQuantity
		if line.Quantity > remainingQuantity {
			rejectedReasons = append(rejectedReasons, fmt.Sprintf(
				"product %s requests %d units but only %d remain returnable",
				product.ID,
				line.Quantity,
				remainingQuantity,
			))
			continue
		}

		acceptedProducts = append(acceptedProducts, product.ID)
	}

	if len(acceptedProducts) == 0 {
		returnReject(memory, strings.Join(rejectedReasons, "; "))
		return nil
	}

	severity := "return-allowed"
	message := fmt.Sprintf(
		"Products %s are eligible for return; requested by %s",
		strings.Join(acceptedProducts, ", "),
		request.RequestedBy.ID,
	)
	if len(rejectedReasons) > 0 {
		severity = "return-partial"
		message = fmt.Sprintf(
			"Products %s are eligible; rejected lines: %s; requested by %s",
			strings.Join(acceptedProducts, ", "),
			strings.Join(rejectedReasons, "; "),
			request.RequestedBy.ID,
		)
	}

	memory.AddFinding(engine.Finding{RuleName: rule.Name(), Severity: severity, Message: message})
	return nil
}

func returnLines(request engine.ReturnRequestFact) []engine.ReturnLineFact {
	if len(request.Lines) > 0 {
		return request.Lines
	}

	return []engine.ReturnLineFact{{
		ProductID:                  request.ProductID,
		Quantity:                   request.Quantity,
		ShippedQuantity:            request.ShippedQuantity,
		PreviouslyReturnedQuantity: request.PreviouslyReturnedQuantity,
	}}
}

func productForReturn(memory *engine.WorkingMemory, productID string) (engine.ProductFact, bool) {
	for _, product := range memory.Products {
		if product.ID == productID {
			return product, true
		}
	}

	return engine.ProductFact{}, false
}

func returnReject(memory *engine.WorkingMemory, message string) {
	memory.AddFinding(engine.Finding{
		RuleName: "Return Policy Rule",
		Severity: "return-rejected",
		Message:  message,
	})
}
