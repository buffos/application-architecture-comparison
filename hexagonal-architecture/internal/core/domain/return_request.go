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
var ErrReturnLineInvalid = errors.New("return line is invalid")
var ErrReturnQuantityExceedsRemaining = errors.New("return quantity exceeds remaining returnable quantity")
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

type ReturnLineRequest struct {
	SKU      string
	Quantity int
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

func NewReturnRequest(order Order, reason, requestedBy string, requestedAt time.Time, requested ...ReturnLineRequest) (ReturnRequest, error) {
	if order.Status != OrderStatusShipped && order.Status != OrderStatusPartiallyShipped {
		return ReturnRequest{}, ErrReturnNotEligible
	}
	if requestedBy == "" {
		return ReturnRequest{}, ErrActorRequired
	}

	id := atomic.AddUint64(&returnSequence, 1)
	orderLineBySKU := make(map[string]OrderLine, len(order.Lines))
	for _, line := range order.Lines {
		orderLineBySKU[line.SKU] = line
		if line.ProductCategory == "Clearance" {
			return ReturnRequest{}, ErrReturnNotEligible
		}
	}

	lines := make([]ReturnLine, 0, len(order.Lines))
	if len(requested) == 0 {
		for _, line := range order.Lines {
			if line.RemainingReturnableQuantity() <= 0 {
				continue
			}
			lines = append(lines, ReturnLine{
				SKU:              line.SKU,
				ProductName:      line.ProductName,
				ProductCategory:  line.ProductCategory,
				Quantity:         line.RemainingReturnableQuantity(),
				LineTotal:        line.AdjustedUnitPrice * line.RemainingReturnableQuantity(),
				ReturnWindowDays: line.ReturnWindowDays,
			})
		}
	} else {
		for _, requestedLine := range requested {
			if requestedLine.Quantity <= 0 {
				return ReturnRequest{}, ErrReturnLineInvalid
			}
			line, ok := orderLineBySKU[requestedLine.SKU]
			if !ok {
				return ReturnRequest{}, ErrReturnLineInvalid
			}
			if requestedLine.Quantity > line.RemainingReturnableQuantity() {
				return ReturnRequest{}, ErrReturnQuantityExceedsRemaining
			}
			lines = append(lines, ReturnLine{
				SKU:              line.SKU,
				ProductName:      line.ProductName,
				ProductCategory:  line.ProductCategory,
				Quantity:         requestedLine.Quantity,
				LineTotal:        line.AdjustedUnitPrice * requestedLine.Quantity,
				ReturnWindowDays: line.ReturnWindowDays,
			})
		}
	}

	if len(lines) == 0 {
		return ReturnRequest{}, ErrReturnNotEligible
	}

	linesCopy := make([]ReturnLine, len(lines))
	copy(linesCopy, lines)
	return ReturnRequest{
		ID:          fmt.Sprintf("ret-%03d", id),
		OrderID:     order.ID,
		RequestedBy: requestedBy,
		Reason:      reason,
		Status:      ReturnStatusRequested,
		RequestedAt: requestedAt,
		ShippedAt:   order.ShippedAt,
		Lines:       linesCopy,
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
