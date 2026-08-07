package returns

import (
	"errors"
	"testing"
)

func TestReturnRequestUsesExplicitReviewBoundary(t *testing.T) {
	accepted := shippedRequestForReview(t)
	if err := accepted.Accept(); err != nil {
		t.Fatal(err)
	}
	if accepted.Status() != ReturnStatusAccepted {
		t.Fatalf("status = %s, want %s", accepted.Status(), ReturnStatusAccepted)
	}
	if err := accepted.Reject(); !errors.Is(err, ErrReturnNotReviewable) {
		t.Fatalf("reviewed request rejected with %v", err)
	}

	rejected := shippedRequestForReview(t)
	if err := rejected.Reject(); err != nil {
		t.Fatal(err)
	}
	if rejected.Status() != ReturnStatusRejected {
		t.Fatalf("status = %s, want %s", rejected.Status(), ReturnStatusRejected)
	}
}

func TestReturnRequestRejectsUnknownReviewDecision(t *testing.T) {
	request := shippedRequestForReview(t)
	if err := request.Review(ReviewDecision("Escalate")); !errors.Is(err, ErrReviewDecisionInvalid) {
		t.Fatalf("unknown decision returned %v", err)
	}
	if request.Status() != ReturnStatusRequested {
		t.Fatal("invalid decision changed request status")
	}
}

func shippedRequestForReview(t *testing.T) ReturnRequest {
	t.Helper()
	returnRequest, err := NewReturnRequestFromShippedOrder("return-review", shippedOrderForReturn(t), "damaged")
	if err != nil {
		t.Fatal(err)
	}
	return returnRequest
}
