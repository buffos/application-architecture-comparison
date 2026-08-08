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
	LineCount                  int
	Partial                    bool
	Action                     engine.ReturnAction
	Outcome                    engine.DecisionOutcome
	Reason                     string
}

func ProjectReturn(
	memory *engine.WorkingMemory,
	decision engine.PolicyDecision,
) ReturnView {
	request := memory.ReturnRequest
	requestedQuantity := request.Quantity
	shippedQuantity := request.ShippedQuantity
	previouslyReturnedQuantity := request.PreviouslyReturnedQuantity
	lineCount := 0
	productID := request.ProductID
	if len(request.Lines) > 0 {
		lineCount = len(request.Lines)
		requestedQuantity = 0
		shippedQuantity = 0
		previouslyReturnedQuantity = 0
		if productID == "" {
			productID = request.Lines[0].ProductID
		}
		for _, line := range request.Lines {
			requestedQuantity += line.Quantity
			shippedQuantity += line.ShippedQuantity
			previouslyReturnedQuantity += line.PreviouslyReturnedQuantity
		}
	} else if request.Requested {
		lineCount = 1
	}
	view := ReturnView{
		OrderID:                    memory.Order.ID,
		ProductID:                  productID,
		Requested:                  request.Requested,
		RequesterID:                request.RequestedBy.ID,
		RequesterRole:              request.RequestedBy.Role,
		RequestedQuantity:          requestedQuantity,
		ShippedQuantity:            shippedQuantity,
		PreviouslyReturnedQuantity: previouslyReturnedQuantity,
		RemainingQuantity:          shippedQuantity - previouslyReturnedQuantity,
		LineCount:                  lineCount,
		Partial:                    decision.ReturnAction == engine.ReturnPartiallyAllowed,
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
