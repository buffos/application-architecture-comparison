package returns

import (
	"errors"
	"fmt"
	"sync/atomic"
	"time"
)

var ErrReturnRequestNotFound = errors.New("return request not found")

var (
	ErrReturnRequestNotReviewable = errors.New("return request is not reviewable")
	ErrActorRequired              = errors.New("actor is required")
)

const (
	ReturnRequestStatusRequested = "Requested"
	ReturnRequestStatusRejected  = "Rejected"
	ReturnRequestStatusRefunded  = "Refunded"
)

var returnRequestSequence uint64

type ReturnRequest struct {
	ID          string
	OrderID     string
	CustomerID  string
	Reason      string
	Status      string
	ShippedAt   time.Time
	RequestedAt time.Time
	RequestedBy string
	ReviewedBy  string
	ProcessedBy string
	ReviewNote  string
	Lines       []ReturnRequestLine
}

type ReturnRequestLine struct {
	ProductSKU       string
	ProductName      string
	ProductCategory  string
	Quantity         int
	UnitPrice        int
	ReturnWindowDays int
}

func NewRefundedReturnRequest(order ReturnableOrder, reason string, requestedAt time.Time, requestedBy string) (ReturnRequest, error) {
	returnRequest, err := NewRequestedReturnRequest(order, reason, requestedAt, requestedBy)
	if err != nil {
		return ReturnRequest{}, err
	}
	returnRequest.Status = ReturnRequestStatusRefunded
	return returnRequest, nil
}

func NewRequestedReturnRequest(order ReturnableOrder, reason string, requestedAt time.Time, requestedBy string) (ReturnRequest, error) {
	if requestedBy == "" {
		return ReturnRequest{}, ErrActorRequired
	}

	id := atomic.AddUint64(&returnRequestSequence, 1)

	lines := make([]ReturnRequestLine, 0, len(order.Lines))
	for _, line := range order.Lines {
		lines = append(lines, ReturnRequestLine{
			ProductSKU:       line.ProductSKU,
			ProductName:      line.ProductName,
			ProductCategory:  line.ProductCategory,
			Quantity:         line.Quantity,
			UnitPrice:        line.UnitPrice,
			ReturnWindowDays: line.ReturnWindowDays,
		})
	}

	return ReturnRequest{
		ID:          fmt.Sprintf("return-%03d", id),
		OrderID:     order.OrderID,
		CustomerID:  order.CustomerID,
		Reason:      reason,
		Status:      ReturnRequestStatusRequested,
		ShippedAt:   order.ShippedAt,
		RequestedAt: requestedAt,
		RequestedBy: requestedBy,
		Lines:       lines,
	}, nil
}

func (r *ReturnRequest) Reject(reviewedBy, reviewNote string) error {
	if r.Status != ReturnRequestStatusRequested {
		return ErrReturnRequestNotReviewable
	}
	if reviewedBy == "" {
		return ErrActorRequired
	}

	r.Status = ReturnRequestStatusRejected
	r.ReviewedBy = reviewedBy
	r.ReviewNote = reviewNote
	return nil
}

func (r *ReturnRequest) Refund(reviewedBy, processedBy, reviewNote string) error {
	if r.Status != ReturnRequestStatusRequested {
		return ErrReturnRequestNotReviewable
	}
	if reviewedBy == "" || processedBy == "" {
		return ErrActorRequired
	}

	r.Status = ReturnRequestStatusRefunded
	r.ReviewedBy = reviewedBy
	r.ProcessedBy = processedBy
	r.ReviewNote = reviewNote
	return nil
}
