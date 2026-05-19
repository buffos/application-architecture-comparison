package domain

import (
	"errors"
	"fmt"
	"sync/atomic"
	"time"
)

const ReturnStatusRequested = "Requested"
const ReturnStatusAccepted = "Accepted"
const ReturnStatusRejected = "Rejected"
const ReturnStatusRefunded = "Refunded"

var returnSequence uint64

var ErrReturnRequestNotFound = errors.New("return request not found")
var ErrReturnNotEligible = errors.New("return is not eligible")
var ErrRefundFailed = errors.New("refund failed")
var ErrReturnRefundNotAllowed = errors.New("return refund is not allowed")
var ErrReturnReviewNotAllowed = errors.New("return review is not allowed")
var ErrActorRequired = errors.New("actor is required")

type ReturnLine struct {
	SKU              string
	ProductName      string
	ProductCategory  string
	Quantity         int
	LineTotal        int
	ReturnWindowDays int
}

type ReturnRequest struct {
	ID          string
	OrderID     string
	RequestedBy string
	ReviewedBy  string
	ProcessedBy string
	ReviewNote  string
	Reason      string
	Status      string
	RequestedAt time.Time
	ShippedAt   time.Time
	Lines       []ReturnLine
}

func NewReturnRequest(order Order, reason, requestedBy string, requestedAt time.Time) (ReturnRequest, error) {
	if order.Status != OrderStatusShipped {
		return ReturnRequest{}, ErrReturnNotEligible
	}
	if requestedBy == "" {
		return ReturnRequest{}, ErrActorRequired
	}

	id := atomic.AddUint64(&returnSequence, 1)
	lines := make([]ReturnLine, 0, len(order.Lines))

	for _, line := range order.Lines {
		if line.ProductCategory == "Clearance" {
			return ReturnRequest{}, ErrReturnNotEligible
		}

		lines = append(lines, ReturnLine{
			SKU:              line.SKU,
			ProductName:      line.ProductName,
			ProductCategory:  line.ProductCategory,
			Quantity:         line.Quantity,
			LineTotal:        line.LineTotal,
			ReturnWindowDays: line.ReturnWindowDays,
		})
	}

	return ReturnRequest{
		ID:          fmt.Sprintf("ret-%03d", id),
		OrderID:     order.ID,
		RequestedBy: requestedBy,
		Reason:      reason,
		Status:      ReturnStatusRequested,
		RequestedAt: requestedAt,
		ShippedAt:   order.ShippedAt,
		Lines:       lines,
	}, nil
}

func (r *ReturnRequest) Accept(reviewedBy string) error {
	if r.Status != ReturnStatusRequested {
		return ErrReturnReviewNotAllowed
	}
	if reviewedBy == "" {
		return ErrActorRequired
	}

	r.ReviewedBy = reviewedBy
	r.ReviewNote = ""
	r.Status = ReturnStatusAccepted
	return nil
}

func (r *ReturnRequest) Reject(reviewedBy, reviewNote string) error {
	if r.Status != ReturnStatusRequested {
		return ErrReturnReviewNotAllowed
	}
	if reviewedBy == "" {
		return ErrActorRequired
	}

	r.ReviewedBy = reviewedBy
	r.ReviewNote = reviewNote
	r.Status = ReturnStatusRejected
	return nil
}

func (r *ReturnRequest) MarkRefunded(processedBy string) error {
	if r.Status != ReturnStatusAccepted {
		return ErrReturnRefundNotAllowed
	}
	if processedBy == "" {
		return ErrActorRequired
	}

	r.ProcessedBy = processedBy
	r.Status = ReturnStatusRefunded
	return nil
}
