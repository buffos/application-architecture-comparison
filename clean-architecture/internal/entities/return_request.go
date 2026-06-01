package entities

import (
	"fmt"
	"sync/atomic"
)

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
		Status:  ReturnRequestStatusRefunded,
	}, nil
}
