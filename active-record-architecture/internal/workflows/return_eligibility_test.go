package workflows

import (
	"testing"

	"active-record-architecture/internal/records"
)

func TestAcceptReturnRejectsClearanceItemBeforeStateChange(t *testing.T) {
	db, request := requestedReturn(t)
	request.Lines[0].ProductCategory = "Clearance"
	if err := request.Save(); err != nil {
		t.Fatalf("ReturnRequest.Save() error = %v", err)
	}

	if _, err := AcceptReturn(db, request.ID); err != records.ErrReturnNotEligible {
		t.Fatalf("clearance acceptance error = %v, want %v", err, records.ErrReturnNotEligible)
	}
	savedRequest, err := records.FindReturnRequest(db, request.ID)
	if err != nil {
		t.Fatalf("FindReturnRequest() error = %v", err)
	}
	if savedRequest.Status != records.ReturnStatusRequested {
		t.Fatalf("request after clearance refusal = %#v", savedRequest)
	}
	stock, err := records.FindStock(db, "sku-001")
	if err != nil {
		t.Fatalf("FindStock() error = %v", err)
	}
	if stock.OnHand != 4 {
		t.Fatalf("stock after clearance refusal = %d, want 4", stock.OnHand)
	}
}

func TestEvaluateReturnEligibilityDoesNotMutateRecords(t *testing.T) {
	db, request := requestedReturn(t)
	order, err := records.FindOrder(db, request.OrderID)
	if err != nil {
		t.Fatalf("FindOrder() error = %v", err)
	}
	decision := request.EvaluateEligibility(order)
	if !decision.Eligible || len(decision.Reasons) != 0 {
		t.Fatalf("normal eligibility decision = %#v", decision)
	}

	request.Lines[0].ProductCategory = "Clearance"
	decision = request.EvaluateEligibility(order)
	if decision.Eligible || len(decision.Reasons) != 1 || decision.Reasons[0] != records.ReturnEligibilityReasonClearance {
		t.Fatalf("clearance eligibility decision = %#v", decision)
	}
	if request.Status != records.ReturnStatusRequested || order.Lines[0].ReturnedQuantity != 0 {
		t.Fatalf("records changed during evaluation: request=%#v orderLine=%#v", request, order.Lines[0])
	}
}

func TestNormalReturnCanStillBeAcceptedBeforeCompletion(t *testing.T) {
	db, request := requestedReturn(t)
	accepted, err := AcceptReturn(db, request.ID)
	if err != nil {
		t.Fatalf("AcceptReturn() error = %v", err)
	}
	if accepted.Status != records.ReturnStatusAccepted {
		t.Fatalf("accepted status = %q, want %q", accepted.Status, records.ReturnStatusAccepted)
	}
}
