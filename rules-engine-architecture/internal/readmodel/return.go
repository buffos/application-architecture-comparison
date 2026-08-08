package readmodel

import "rules-engine-architecture/internal/engine"

// ReturnView is a stable read shape for a return request. It contains no
// slices or pointers into WorkingMemory.
type ReturnView struct {
	OrderID                    string
	ProductID                  string
	Requested                  bool
	RequesterID                string
	RequesterRole              string
	RequestedQuantity          int
	ShippedQuantity            int
	PreviouslyReturnedQuantity int
	RemainingQuantity          int
	Action                     engine.ReturnAction
	Outcome                    engine.DecisionOutcome
	Reason                     string
}

func ProjectReturn(
	memory *engine.WorkingMemory,
	decision engine.PolicyDecision,
) ReturnView {
	request := memory.ReturnRequest
	view := ReturnView{
		OrderID:                    memory.Order.ID,
		ProductID:                  request.ProductID,
		Requested:                  request.Requested,
		RequesterID:                request.RequestedBy.ID,
		RequesterRole:              request.RequestedBy.Role,
		RequestedQuantity:          request.Quantity,
		ShippedQuantity:            request.ShippedQuantity,
		PreviouslyReturnedQuantity: request.PreviouslyReturnedQuantity,
		RemainingQuantity:          request.ShippedQuantity - request.PreviouslyReturnedQuantity,
		Action:                     decision.ReturnAction,
		Outcome:                    decision.Outcome,
	}

	for _, reason := range decision.Reasons {
		if reason.RuleName == "Return Policy Rule" {
			view.Reason = reason.Message
			break
		}
	}

	return view
}
