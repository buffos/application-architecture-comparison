package domain

import (
	"errors"
	"fmt"
	"sync/atomic"
	"time"
)

var ErrReturnRequestNotFound = errors.New("return request not found")
var ErrReturnRequestNotAcceptable = errors.New("return request is not acceptable")
var ErrReturnRequestNotRejectable = errors.New("return request is not rejectable")
var ErrActorRequired = errors.New("actor is required")

const ReturnRequestStatusRequested = "Requested"
const ReturnRequestStatusAccepted = "Accepted"
const ReturnRequestStatusRejected = "Rejected"
const ReturnRequestStatusRefunded = "Refunded"

var returnRequestSequence uint64

type ReturnRequest struct {
	ID          string
	OrderID     string
	Reason      string
	Status      string
	RequestedAt time.Time
	RequestedBy string
	ReviewedBy  string
	ProcessedBy string
	ReviewNote  string
}

func NewReturnRequest(orderID string, reason string, requestedAt time.Time, requestedBy string) (ReturnRequest, error) {
	if requestedBy == "" {
		return ReturnRequest{}, ErrActorRequired
	}

	id := atomic.AddUint64(&returnRequestSequence, 1)

	return ReturnRequest{
		ID:          fmt.Sprintf("return-%03d", id),
		OrderID:     orderID,
		Reason:      reason,
		Status:      ReturnRequestStatusRequested,
		RequestedAt: requestedAt,
		RequestedBy: requestedBy,
	}, nil
}

func (r *ReturnRequest) Accept(reviewedBy string, reviewNote string) error {
	if r.Status != ReturnRequestStatusRequested {
		return ErrReturnRequestNotAcceptable
	}

	if reviewedBy == "" {
		return ErrActorRequired
	}

	r.Status = ReturnRequestStatusAccepted
	r.ReviewedBy = reviewedBy
	r.ReviewNote = reviewNote
	return nil
}

func (r *ReturnRequest) Reject(reviewedBy string, reviewNote string) error {
	if r.Status != ReturnRequestStatusRequested {
		return ErrReturnRequestNotRejectable
	}

	if reviewedBy == "" {
		return ErrActorRequired
	}

	r.Status = ReturnRequestStatusRejected
	r.ReviewedBy = reviewedBy
	r.ReviewNote = reviewNote
	return nil
}

func (r *ReturnRequest) MarkRefunded(processedBy string) error {
	if r.Status != ReturnRequestStatusAccepted {
		return ErrReturnRequestNotAcceptable
	}

	if processedBy == "" {
		return ErrActorRequired
	}

	r.Status = ReturnRequestStatusRefunded
	r.ProcessedBy = processedBy
	return nil
}
