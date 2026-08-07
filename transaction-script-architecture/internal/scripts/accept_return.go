package scripts

import (
	"errors"
	"fmt"

	"transaction-script-architecture/internal/data"
)

var (
	ErrReturnNotAcceptable = errors.New("return is not requested")
	ErrReturnOrderMissing  = errors.New("return order not found")
	ErrReturnStockMissing  = errors.New("return stock record not found")
)

// AcceptReturn completes the initial return workflow: it updates the order,
// restocks inventory, and creates a completed refund.
func AcceptReturn(store *data.Store, returnID string) (data.ReturnRequest, error) {
	if store == nil {
		return data.ReturnRequest{}, ErrStoreRequired
	}

	if returnID == "" {
		return data.ReturnRequest{}, ErrReturnIDRequired
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

	for _, returnLine := range request.Lines {
		stock, ok := store.Stocks[returnLine.SKU]
		if !ok {
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
			_ = stock
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
	store.Orders[order.ID] = order
	store.Refunds[refund.ID] = refund
	store.Returns[request.ID] = request

	return request, nil
}
