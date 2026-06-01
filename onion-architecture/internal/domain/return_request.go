package domain

import (
	"errors"
	"fmt"
	"sync/atomic"
)

var ErrReturnRequestNotFound = errors.New("return request not found")

const ReturnRequestStatusRequested = "Requested"
const ReturnRequestStatusRefunded = "Refunded"

var returnRequestSequence uint64

type ReturnRequest struct {
	ID      string
	OrderID string
	Reason  string
	Status  string
}

func NewReturnRequest(orderID string, reason string) (ReturnRequest, error) {
	id := atomic.AddUint64(&returnRequestSequence, 1)

	return ReturnRequest{
		ID:      fmt.Sprintf("return-%03d", id),
		OrderID: orderID,
		Reason:  reason,
		Status:  ReturnRequestStatusRequested,
	}, nil
}
