package scripts

import (
	"errors"
	"fmt"
	"time"

	"transaction-script-architecture/internal/data"
)

var (
	ErrOrderNotReturnable = errors.New("order is not shipped")
	ErrReturnLinesInvalid = errors.New("return lines are invalid")
	ErrReturnIDRequired   = errors.New("return id is required")
	ErrReturnNotFound     = errors.New("return request not found")
	ErrActorRequired      = errors.New("return actor is required")
)

// RequestReturn creates a return request after checking the order's shipped
// and already-returned quantities. It does not refund or restock anything.
func RequestReturn(store *data.Store, orderID string, lines []data.ReturnLine, reason string, requestedBy string) (data.ReturnRequest, error) {
	return RequestReturnAt(store, orderID, lines, reason, requestedBy, time.Now())
}

// RequestReturnAt is the deterministic form used by tests and demonstrations.
func RequestReturnAt(store *data.Store, orderID string, lines []data.ReturnLine, reason string, requestedBy string, requestedAt time.Time) (data.ReturnRequest, error) {
	if store == nil {
		return data.ReturnRequest{}, ErrStoreRequired
	}

	if orderID == "" {
		return data.ReturnRequest{}, ErrOrderIDRequired
	}

	if requestedBy == "" {
		return data.ReturnRequest{}, ErrActorRequired
	}

	order, ok := store.Orders[orderID]
	if !ok {
		return data.ReturnRequest{}, ErrOrderNotFound
	}

	if order.Status != data.OrderStatusShipped && order.Status != data.OrderStatusPartiallyShipped {
		return data.ReturnRequest{}, ErrOrderNotReturnable
	}

	requestLines := lines
	if len(requestLines) == 0 {
		requestLines = make([]data.ReturnLine, 0, len(order.Lines))
		for _, orderLine := range order.Lines {
			remaining := orderLine.ShippedQuantity - orderLine.ReturnedQuantity
			if remaining <= 0 {
				continue
			}
			requestLines = append(requestLines, data.ReturnLine{
				OrderLineID:      orderLine.ID,
				SKU:              orderLine.SKU,
				ProductCategory:  orderLine.ProductCategory,
				Quantity:         remaining,
				UnitPrice:        orderLine.UnitPrice,
				ReturnWindowDays: orderLine.ReturnWindowDays,
			})
		}
	}

	if len(requestLines) == 0 {
		return data.ReturnRequest{}, ErrReturnLinesInvalid
	}

	refundAmount := 0
	for _, requestedLine := range requestLines {
		if requestedLine.Quantity <= 0 {
			return data.ReturnRequest{}, ErrReturnLinesInvalid
		}

		matched := false
		for _, orderLine := range order.Lines {
			if orderLine.ID != requestedLine.OrderLineID {
				continue
			}

			remaining := orderLine.ShippedQuantity - orderLine.ReturnedQuantity
			if requestedLine.Quantity > remaining {
				return data.ReturnRequest{}, ErrReturnLinesInvalid
			}

			matched = true
			refundAmount += requestedLine.Quantity * orderLine.UnitPrice
			break
		}

		if !matched {
			return data.ReturnRequest{}, ErrReturnLinesInvalid
		}
	}

	store.NextReturnNumber++
	request := data.ReturnRequest{
		ID:           fmt.Sprintf("return-%03d", store.NextReturnNumber),
		OrderID:      order.ID,
		Status:       data.ReturnStatusRequested,
		Reason:       reason,
		Lines:        requestLines,
		RefundStatus: data.RefundStatusNotStarted,
		RefundAmount: refundAmount,
		RequestedAt:  requestedAt,
		RequestedBy:  requestedBy,
	}
	store.Returns[request.ID] = request

	return request, nil
}
