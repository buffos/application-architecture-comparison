package entities

import (
	"fmt"
	"sync/atomic"
)

const ReturnRequestStatusRequested = "Requested"
const ReturnRequestStatusAccepted = "Accepted"
const ReturnRequestStatusRejected = "Rejected"
const ReturnRequestStatusRefunded = "Refunded"

var returnSequence uint64

var ErrOrderNotReturnable = ErrQuoteCannotTransition

type ReturnRequest struct {
	ID      string
	OrderID string
	Reason  string
	Status  string
}

func NewReturnRequestFromShippedOrder(order Order, reason string) (ReturnRequest, error) {
	if order.Status != OrderStatusShipped {
		return ReturnRequest{}, ErrOrderNotReturnable
	}

	id := atomic.AddUint64(&returnSequence, 1)

	return ReturnRequest{
		ID:      fmt.Sprintf("return-%03d", id),
		OrderID: order.ID,
		Reason:  reason,
		Status:  ReturnRequestStatusRequested,
	}, nil
}

func (r *ReturnRequest) Accept() error {
	if r.Status != ReturnRequestStatusRequested {
		return ErrQuoteCannotTransition
	}

	r.Status = ReturnRequestStatusAccepted
	return nil
}

func (r *ReturnRequest) Reject() error {
	if r.Status != ReturnRequestStatusRequested {
		return ErrQuoteCannotTransition
	}

	r.Status = ReturnRequestStatusRejected
	return nil
}

func (r *ReturnRequest) MarkRefunded() error {
	if r.Status != ReturnRequestStatusAccepted {
		return ErrQuoteCannotTransition
	}

	r.Status = ReturnRequestStatusRefunded
	return nil
}
