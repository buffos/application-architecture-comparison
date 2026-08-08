package rules

import (
	"fmt"

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

	product, found := productForReturn(memory, request.ProductID)
	if !found {
		returnReject(memory, fmt.Sprintf(
			"Return request references unknown product %s",
			request.ProductID,
		))
		return nil
	}

	if product.Category == "Clearance" {
		returnReject(memory, fmt.Sprintf("Product %s is clearance and cannot be returned", product.ID))
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

	if request.Quantity <= 0 {
		returnReject(memory, "Return quantity must be greater than zero")
		return nil
	}

	remainingQuantity := request.ShippedQuantity - request.PreviouslyReturnedQuantity
	if request.Quantity > remainingQuantity {
		returnReject(memory, fmt.Sprintf(
			"Requested return quantity %d exceeds remaining returnable quantity %d",
			request.Quantity,
			remainingQuantity,
		))
		return nil
	}

	memory.AddFinding(engine.Finding{
		RuleName: rule.Name(),
		Severity: "return-allowed",
		Message: fmt.Sprintf(
			"Product %s is eligible for return; requested by %s",
			product.ID,
			request.RequestedBy.ID,
		),
	})
	return nil
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
