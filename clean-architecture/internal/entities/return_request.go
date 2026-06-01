package entities

import (
	"errors"
	"fmt"
	"sync/atomic"
	"time"
)

const ReturnRequestStatusRequested = "Requested"
const ReturnRequestStatusAccepted = "Accepted"
const ReturnRequestStatusRejected = "Rejected"
const ReturnRequestStatusRefunded = "Refunded"

var returnSequence uint64

var ErrOrderNotReturnable = ErrQuoteCannotTransition
var ErrActorRequired = errors.New("actor is required")

type ReturnRequestLine struct {
	SKU         string
	ProductName string
	Quantity    int
}

type ReturnRequest struct {
	ID          string
	OrderID     string
	Reason      string
	Lines       []ReturnRequestLine
	Status      string
	RequestedAt time.Time
	RequestedBy string
	ReviewedBy  string
	ProcessedBy string
	ReviewNote  string
}

func NewReturnRequestFromShippedOrder(order Order, reason string, lines []ReturnRequestLine, requestedAt time.Time, requestedBy string) (ReturnRequest, error) {
	if order.Status != OrderStatusShipped && order.Status != OrderStatusPartiallyShipped {
		return ReturnRequest{}, ErrOrderNotReturnable
	}
	if requestedBy == "" {
		return ReturnRequest{}, ErrActorRequired
	}

	requestLines := lines
	if len(requestLines) == 0 {
		requestLines = make([]ReturnRequestLine, 0, len(order.Lines))
		for _, line := range order.Lines {
			returnable := line.ShippedQuantity - line.ReturnedQuantity
			if returnable <= 0 {
				continue
			}

			requestLines = append(requestLines, ReturnRequestLine{
				SKU:         line.SKU,
				ProductName: line.ProductName,
				Quantity:    returnable,
			})
		}
	}

	if len(requestLines) == 0 {
		return ReturnRequest{}, ErrOrderNotReturnable
	}

	for _, requestLine := range requestLines {
		if requestLine.Quantity <= 0 {
			return ReturnRequest{}, ErrQuoteCannotTransition
		}

		matched := false
		for _, orderLine := range order.Lines {
			if orderLine.SKU != requestLine.SKU {
				continue
			}

			returnable := orderLine.ShippedQuantity - orderLine.ReturnedQuantity
			if requestLine.Quantity > returnable {
				return ReturnRequest{}, ErrQuoteCannotTransition
			}

			matched = true
			break
		}

		if !matched {
			return ReturnRequest{}, ErrQuoteCannotTransition
		}
	}

	id := atomic.AddUint64(&returnSequence, 1)

	return ReturnRequest{
		ID:          fmt.Sprintf("return-%03d", id),
		OrderID:     order.ID,
		Reason:      reason,
		Lines:       requestLines,
		Status:      ReturnRequestStatusRequested,
		RequestedAt: requestedAt,
		RequestedBy: requestedBy,
	}, nil
}

func (r *ReturnRequest) Accept(reviewedBy string) error {
	if r.Status != ReturnRequestStatusRequested {
		return ErrQuoteCannotTransition
	}
	if reviewedBy == "" {
		return ErrActorRequired
	}

	r.Status = ReturnRequestStatusAccepted
	r.ReviewedBy = reviewedBy
	return nil
}

func (r *ReturnRequest) Reject(reviewedBy string, reviewNote string) error {
	if r.Status != ReturnRequestStatusRequested {
		return ErrQuoteCannotTransition
	}
	if reviewedBy == "" {
		return ErrActorRequired
	}

	r.Status = ReturnRequestStatusRejected
	r.ReviewedBy = reviewedBy
	r.ReviewNote = reviewNote
	return nil
}

func (r *ReturnRequest) MarkRefunded(processedBy string) error {
	if r.Status != ReturnRequestStatusAccepted {
		return ErrQuoteCannotTransition
	}
	if processedBy == "" {
		return ErrActorRequired
	}

	r.Status = ReturnRequestStatusRefunded
	r.ProcessedBy = processedBy
	return nil
}
