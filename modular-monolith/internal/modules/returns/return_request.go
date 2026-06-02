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
	ErrReturnQuantityInvalid      = errors.New("return quantity is invalid")
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

type RequestedReturnLine struct {
	ProductSKU string
	Quantity   int
}

func NewRefundedReturnRequest(order ReturnableOrder, reason string, requestedAt time.Time, requestedBy string) (ReturnRequest, error) {
	returnRequest, err := NewRequestedReturnRequest(order, nil, reason, requestedAt, requestedBy)
	if err != nil {
		return ReturnRequest{}, err
	}
	returnRequest.Status = ReturnRequestStatusRefunded
	return returnRequest, nil
}

func NewRequestedReturnRequest(order ReturnableOrder, requestedLines []RequestedReturnLine, reason string, requestedAt time.Time, requestedBy string) (ReturnRequest, error) {
	if requestedBy == "" {
		return ReturnRequest{}, ErrActorRequired
	}

	id := atomic.AddUint64(&returnRequestSequence, 1)

	if len(requestedLines) == 0 {
		for _, line := range order.Lines {
			remaining := line.ShippedQuantity
			if remaining > 0 {
				requestedLines = append(requestedLines, RequestedReturnLine{
					ProductSKU: line.ProductSKU,
					Quantity:   remaining,
				})
			}
		}
	}

	lines := make([]ReturnRequestLine, 0, len(requestedLines))
	for _, requestedLine := range requestedLines {
		if requestedLine.Quantity <= 0 {
			return ReturnRequest{}, ErrReturnQuantityInvalid
		}

		matched := false
		for _, line := range order.Lines {
			if line.ProductSKU != requestedLine.ProductSKU {
				continue
			}

			if requestedLine.Quantity > line.ShippedQuantity {
				return ReturnRequest{}, ErrReturnQuantityInvalid
			}

			lines = append(lines, ReturnRequestLine{
				ProductSKU:       line.ProductSKU,
				ProductName:      line.ProductName,
				ProductCategory:  line.ProductCategory,
				Quantity:         requestedLine.Quantity,
				UnitPrice:        line.UnitPrice,
				ReturnWindowDays: line.ReturnWindowDays,
			})
			matched = true
			break
		}

		if !matched {
			return ReturnRequest{}, ErrReturnQuantityInvalid
		}
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
