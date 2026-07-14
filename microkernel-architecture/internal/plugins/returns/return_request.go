package returns

import (
	"errors"
	"fmt"
	"sync/atomic"
	"time"
)

var ErrReturnRequestNotFound = errors.New("return request not found")
var ErrReturnRequestNotReviewable = errors.New("return request is not reviewable")

const ReturnRequestStatusRequested = "Requested"
const ReturnRequestStatusRefunded = "Refunded"
const ReturnRequestStatusRejected = "Rejected"

var returnSequence uint64

type ReturnRequest struct {
	ID          string
	OrderID     string
	CustomerID  string
	Reason      string
	ShippedAt   time.Time
	RequestedAt time.Time
	Status      string
	Lines       []ReturnLine
}

type ReturnLine struct {
	ProductSKU       string
	Quantity         int
	UnitPrice        int
	ReturnWindowDays int
}

func NewReturnRequest(orderID string, customerID string, reason string, shippedAt time.Time, requestedAt time.Time, lines []ReturnLine) ReturnRequest {
	id := atomic.AddUint64(&returnSequence, 1)

	return ReturnRequest{
		ID:          fmt.Sprintf("return-%03d", id),
		OrderID:     orderID,
		CustomerID:  customerID,
		Reason:      reason,
		ShippedAt:   shippedAt,
		RequestedAt: requestedAt,
		Status:      ReturnRequestStatusRequested,
		Lines:       lines,
	}
}

func (r ReturnRequest) TotalAmount() int {
	total := 0
	for _, line := range r.Lines {
		total += line.Quantity * line.UnitPrice
	}

	return total
}

func (r *ReturnRequest) Accept() error {
	if r.Status != ReturnRequestStatusRequested {
		return ErrReturnRequestNotReviewable
	}

	r.Status = ReturnRequestStatusRefunded
	return nil
}

func (r *ReturnRequest) Reject() error {
	if r.Status != ReturnRequestStatusRequested {
		return ErrReturnRequestNotReviewable
	}

	r.Status = ReturnRequestStatusRejected
	return nil
}
