package scripts

import (
	"errors"
	"fmt"
	"time"

	"transaction-script-architecture/internal/data"
)

var (
	ErrReturnNotAcceptable = errors.New("return is not requested")
	ErrReturnNotRejectable = errors.New("return is not requested")
	ErrReturnNotRefundable = errors.New("return must be accepted before refund")
	ErrReturnNotEligible   = errors.New("return is not eligible")
	ErrReturnOrderMissing  = errors.New("return order not found")
	ErrReturnStockMissing  = errors.New("return stock record not found")
)

// AcceptReturn records the review decision. It intentionally performs no
// financial or inventory side effect until CompleteRefund runs.
func AcceptReturn(store *data.Store, returnID string, reviewedBy string) (data.ReturnRequest, error) {
	return AcceptReturnAt(store, returnID, time.Now(), reviewedBy)
}

// AcceptReturnAt is the deterministic form used by tests and demonstrations.
func AcceptReturnAt(store *data.Store, returnID string, now time.Time, reviewedBy string) (data.ReturnRequest, error) {
	if store == nil {
		return data.ReturnRequest{}, ErrStoreRequired
	}

	if returnID == "" {
		return data.ReturnRequest{}, ErrReturnIDRequired
	}

	if reviewedBy == "" {
		return data.ReturnRequest{}, ErrActorRequired
	}

	request, ok := store.Returns[returnID]
	if !ok {
		return data.ReturnRequest{}, ErrReturnNotFound
	}

	if request.Status != data.ReturnStatusRequested {
		return data.ReturnRequest{}, ErrReturnNotAcceptable
	}

	order, ok := store.Orders[request.OrderID]
	if !ok {
		return data.ReturnRequest{}, ErrReturnOrderMissing
	}
	if decision := EvaluateReturnEligibilityAt(order, request, now); !decision.Eligible {
		return data.ReturnRequest{}, ErrReturnNotEligible
	}

	request.Status = data.ReturnStatusAccepted
	request.ReviewedBy = reviewedBy
	store.Returns[request.ID] = request

	return request, nil
}

// RejectReturn records a negative review decision and blocks later refund
// processing.
func RejectReturn(store *data.Store, returnID string, reviewedBy string, reviewNote string) (data.ReturnRequest, error) {
	if store == nil {
		return data.ReturnRequest{}, ErrStoreRequired
	}

	if returnID == "" {
		return data.ReturnRequest{}, ErrReturnIDRequired
	}

	if reviewedBy == "" {
		return data.ReturnRequest{}, ErrActorRequired
	}

	request, ok := store.Returns[returnID]
	if !ok {
		return data.ReturnRequest{}, ErrReturnNotFound
	}

	if request.Status != data.ReturnStatusRequested {
		return data.ReturnRequest{}, ErrReturnNotRejectable
	}

	request.Status = data.ReturnStatusRejected
	request.ReviewedBy = reviewedBy
	request.ReviewNote = reviewNote
	store.Returns[request.ID] = request

	return request, nil
}

// CompleteRefund applies the accepted return's order, inventory, and refund
// writes as one procedural transaction.
func CompleteRefund(store *data.Store, returnID string, processedBy string) (data.ReturnRequest, error) {
	if store == nil {
		return data.ReturnRequest{}, ErrStoreRequired
	}

	if returnID == "" {
		return data.ReturnRequest{}, ErrReturnIDRequired
	}

	if processedBy == "" {
		return data.ReturnRequest{}, ErrActorRequired
	}

	request, ok := store.Returns[returnID]
	if !ok {
		return data.ReturnRequest{}, ErrReturnNotFound
	}

	if request.Status != data.ReturnStatusAccepted {
		return data.ReturnRequest{}, ErrReturnNotRefundable
	}

	order, ok := store.Orders[request.OrderID]
	if !ok {
		return data.ReturnRequest{}, ErrReturnOrderMissing
	}

	for _, returnLine := range request.Lines {
		if _, ok := store.Stocks[returnLine.SKU]; !ok {
			return data.ReturnRequest{}, ErrReturnStockMissing
		}

		matched := false
		for _, orderLine := range order.Lines {
			if orderLine.ID != returnLine.OrderLineID {
				continue
			}

			remaining := orderLine.ShippedQuantity - orderLine.ReturnedQuantity
			if returnLine.Quantity > remaining {
				return data.ReturnRequest{}, ErrReturnLinesInvalid
			}
			matched = true
			break
		}
		if !matched {
			return data.ReturnRequest{}, ErrReturnLinesInvalid
		}
	}

	store.NextRefundNumber++
	refund := data.Refund{
		ID:              fmt.Sprintf("refund-%03d", store.NextRefundNumber),
		ReturnRequestID: request.ID,
		OrderID:         order.ID,
		Amount:          request.RefundAmount,
		Status:          data.RefundStatusCompleted,
		ProcessedBy:     processedBy,
	}

	for _, returnLine := range request.Lines {
		for index := range order.Lines {
			if order.Lines[index].ID != returnLine.OrderLineID {
				continue
			}

			order.Lines[index].ReturnedQuantity += returnLine.Quantity
			stock := store.Stocks[returnLine.SKU]
			stock.OnHand += returnLine.Quantity
			store.Stocks[returnLine.SKU] = stock
			break
		}
	}

	request.Status = data.ReturnStatusRefunded
	request.RefundID = refund.ID
	request.RefundStatus = refund.Status
	request.ProcessedBy = processedBy
	store.Orders[order.ID] = order
	store.Refunds[refund.ID] = refund
	store.Returns[request.ID] = request

	return request, nil
}
