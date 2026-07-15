package returns

import (
	"errors"
	"fmt"
	"sync/atomic"
	"time"
)

var ErrReturnRequestNotFound = errors.New("return request not found")
var ErrReturnRequestNotReviewable = errors.New("return request is not reviewable")
var ErrActorRequired = errors.New("actor is required")

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
	RequestedBy string
	ReviewedBy  string
	ProcessedBy string
	ReviewNote  string
	Status      string
	Lines       []ReturnLine
}

type ReturnLine struct {
	ProductSKU       string
	ProductCategory  string
	Quantity         int
	UnitPrice        int
	ReturnWindowDays int
}

func NewReturnRequest(orderID string, customerID string, reason string, shippedAt time.Time, requestedAt time.Time, requestedBy string, lines []ReturnLine) (ReturnRequest, error) {
	if requestedBy == "" {
		return ReturnRequest{}, ErrActorRequired
	}

	id := atomic.AddUint64(&returnSequence, 1)

	return ReturnRequest{
		ID:          fmt.Sprintf("return-%03d", id),
		OrderID:     orderID,
		CustomerID:  customerID,
		Reason:      reason,
		ShippedAt:   shippedAt,
		RequestedAt: requestedAt,
		RequestedBy: requestedBy,
		Status:      ReturnRequestStatusRequested,
		Lines:       lines,
	}, nil
}

func (r ReturnRequest) TotalAmount() int {
	total := 0
	for _, line := range r.Lines {
		total += line.Quantity * line.UnitPrice
	}

	return total
}

func (r *ReturnRequest) Accept(reviewedBy string, processedBy string, reviewNote string) error {
	if r.Status != ReturnRequestStatusRequested {
		return ErrReturnRequestNotReviewable
	}

	if reviewedBy == "" || processedBy == "" {
		return ErrActorRequired
	}

	r.ReviewedBy = reviewedBy
	r.ProcessedBy = processedBy
	r.ReviewNote = reviewNote
	r.Status = ReturnRequestStatusRefunded
	return nil
}

func (r *ReturnRequest) Reject(reviewedBy string, reviewNote string) error {
	if r.Status != ReturnRequestStatusRequested {
		return ErrReturnRequestNotReviewable
	}

	if reviewedBy == "" {
		return ErrActorRequired
	}

	r.ReviewedBy = reviewedBy
	r.ReviewNote = reviewNote
	r.Status = ReturnRequestStatusRejected
	return nil
}
