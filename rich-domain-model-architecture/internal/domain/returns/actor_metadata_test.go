package returns

import (
	"errors"
	"testing"
)

func TestReturnRequestRecordsActorMetadata(t *testing.T) {
	request := shippedRequestForReview(t)
	if err := request.AssignRequester("customer-001"); err != nil {
		t.Fatal(err)
	}
	if err := request.ReviewBy(ReviewDecisionAccept, "reviewer-001"); err != nil {
		t.Fatal(err)
	}
	if err := request.ProcessBy("processor-001"); err != nil {
		t.Fatal(err)
	}
	if request.RequestedBy() != "customer-001" || request.ReviewedBy() != "reviewer-001" || request.ProcessedBy() != "processor-001" {
		t.Fatalf("unexpected metadata: requested=%s reviewed=%s processed=%s", request.RequestedBy(), request.ReviewedBy(), request.ProcessedBy())
	}
}

func TestReturnRequestRejectsMissingActorsAndInvalidProcessing(t *testing.T) {
	request := shippedRequestForReview(t)
	if err := request.AssignRequester(""); !errors.Is(err, ErrActorRequired) {
		t.Fatalf("missing requester returned %v", err)
	}
	if err := request.ReviewBy(ReviewDecisionAccept, ""); !errors.Is(err, ErrActorRequired) {
		t.Fatalf("missing reviewer returned %v", err)
	}
	if err := request.ProcessBy("processor-001"); !errors.Is(err, ErrReturnNotReviewable) {
		t.Fatalf("processing requested return returned %v", err)
	}
	if err := request.AssignRequester("customer-001"); err != nil {
		t.Fatal(err)
	}
	if err := request.ReviewBy(ReviewDecisionReject, "reviewer-001"); err != nil {
		t.Fatal(err)
	}
	if err := request.ProcessBy("processor-001"); !errors.Is(err, ErrReturnNotReviewable) {
		t.Fatalf("processing rejected return returned %v", err)
	}
}
