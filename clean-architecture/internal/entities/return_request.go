package entities

import (
	"errors"
	"fmt"
	"sync/atomic"
	"time"
)

const ReturnRequestStatusRequested = "Requested"
const ReturnRequestStatusAccepted = "Accepted"
const ReturnRequestStatusRejected = "Rejected"
const ReturnRequestStatusRefunded = "Refunded"

var returnSequence uint64

var ErrOrderNotReturnable = ErrQuoteCannotTransition
var ErrActorRequired = errors.New("actor is required")

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

func NewReturnRequestFromShippedOrder(order Order, reason string, requestedAt time.Time, requestedBy string) (ReturnRequest, error) {
	if order.Status != OrderStatusShipped {
		return ReturnRequest{}, ErrOrderNotReturnable
	}
	if requestedBy == "" {
		return ReturnRequest{}, ErrActorRequired
	}

	id := atomic.AddUint64(&returnSequence, 1)

	return ReturnRequest{
		ID:          fmt.Sprintf("return-%03d", id),
		OrderID:     order.ID,
		Reason:      reason,
		Status:      ReturnRequestStatusRequested,
		RequestedAt: requestedAt,
		RequestedBy: requestedBy,
	}, nil
}

func (r *ReturnRequest) Accept(reviewedBy string) error {
	if r.Status != ReturnRequestStatusRequested {
		return ErrQuoteCannotTransition
	}
	if reviewedBy == "" {
		return ErrActorRequired
	}

	r.Status = ReturnRequestStatusAccepted
	r.ReviewedBy = reviewedBy
	return nil
}

func (r *ReturnRequest) Reject(reviewedBy string, reviewNote string) error {
	if r.Status != ReturnRequestStatusRequested {
		return ErrQuoteCannotTransition
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
		return ErrQuoteCannotTransition
	}
	if processedBy == "" {
		return ErrActorRequired
	}

	r.Status = ReturnRequestStatusRefunded
	r.ProcessedBy = processedBy
	return nil
}
