package domain

import (
	"errors"
	"fmt"
	"sync/atomic"
)

const ReturnStatusRequested = "Requested"
const ReturnStatusAccepted = "Accepted"
const ReturnStatusRejected = "Rejected"

var returnRequestSequence uint64

var ErrReturnRequestNotFound = errors.New("return request not found")
var ErrReturnNotEligible = errors.New("return is not eligible")
var ErrReturnAlreadyAccepted = errors.New("return request is already accepted")

type ReturnLine struct {
	SKU      string
	Quantity int
}

type ReturnRequest struct {
	ID      string
	OrderID string
	Status  string
	Reason  string
	Lines   []ReturnLine
}

func NewReturnRequest(order Order, reason string) (ReturnRequest, error) {
	if order.Status != OrderStatusShipped {
		return ReturnRequest{}, ErrReturnNotEligible
	}

	for _, line := range order.Lines {
		if line.ProductCategory == "Clearance" {
			return ReturnRequest{}, ErrReturnNotEligible
		}
	}

	id := atomic.AddUint64(&returnRequestSequence, 1)
	lines := make([]ReturnLine, 0, len(order.Lines))

	for _, line := range order.Lines {
		lines = append(lines, ReturnLine{
			SKU:      line.SKU,
			Quantity: line.Quantity,
		})
	}

	return ReturnRequest{
		ID:      fmt.Sprintf("ret-%03d", id),
		OrderID: order.ID,
		Status:  ReturnStatusRequested,
		Reason:  reason,
		Lines:   lines,
	}, nil
}

func (r *ReturnRequest) Accept() error {
	if r.Status == ReturnStatusAccepted {
		return ErrReturnAlreadyAccepted
	}

	r.Status = ReturnStatusAccepted
	return nil
}
