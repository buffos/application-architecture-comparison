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

const ReturnRequestStatusRequested = "Requested"
const ReturnRequestStatusAccepted = "Accepted"
const ReturnRequestStatusRejected = "Rejected"
const ReturnRequestStatusRefunded = "Refunded"

var returnRequestSequence uint64

type ReturnRequest struct {
	ID      string
	OrderID string
	Reason  string
	Status  string
	RequestedAt time.Time
}

func NewReturnRequest(orderID string, reason string, requestedAt time.Time) (ReturnRequest, error) {
	id := atomic.AddUint64(&returnRequestSequence, 1)

	return ReturnRequest{
		ID:      fmt.Sprintf("return-%03d", id),
		OrderID: orderID,
		Reason:  reason,
		Status:  ReturnRequestStatusRequested,
		RequestedAt: requestedAt,
	}, nil
}

func (r *ReturnRequest) Accept() error {
	if r.Status != ReturnRequestStatusRequested {
		return ErrReturnRequestNotAcceptable
	}

	r.Status = ReturnRequestStatusAccepted
	return nil
}

func (r *ReturnRequest) Reject() error {
	if r.Status != ReturnRequestStatusRequested {
		return ErrReturnRequestNotRejectable
	}

	r.Status = ReturnRequestStatusRejected
	return nil
}

func (r *ReturnRequest) MarkRefunded() error {
	if r.Status != ReturnRequestStatusAccepted {
		return ErrReturnRequestNotAcceptable
	}

	r.Status = ReturnRequestStatusRefunded
	return nil
}
