package scripts

import (
	"errors"

	"transaction-script-architecture/internal/data"
)

var (
	ErrOrderNotCancellable        = errors.New("order cannot be cancelled")
	ErrCancelledByRequired        = errors.New("cancelling actor is required")
	ErrCancellationReasonRequired = errors.New("cancellation reason is required")
	ErrStockReleaseInvalid        = errors.New("reserved stock cannot be released")
)

// CancelOrder cancels an unshipped order and releases its outstanding
// reservations as one procedural transaction.
func CancelOrder(store *data.Store, orderID string, cancelledBy string, reason string) (data.Order, error) {
	if store == nil {
		return data.Order{}, ErrStoreRequired
	}

	if orderID == "" {
		return data.Order{}, ErrOrderIDRequired
	}

	if cancelledBy == "" {
		return data.Order{}, ErrCancelledByRequired
	}

	if reason == "" {
		return data.Order{}, ErrCancellationReasonRequired
	}

	order, ok := store.Orders[orderID]
	if !ok {
		return data.Order{}, ErrOrderNotFound
	}

	if order.Status == data.OrderStatusShipped || order.Status == data.OrderStatusPartiallyShipped || order.Status == data.OrderStatusCancelled {
		return data.Order{}, ErrOrderNotCancellable
	}

	for _, line := range order.Lines {
		if line.ReservedQuantity == 0 {
			continue
		}

		stock, ok := store.Stocks[line.SKU]
		if !ok || stock.Reserved < line.ReservedQuantity {
			return data.Order{}, ErrStockReleaseInvalid
		}
	}

	for index := range order.Lines {
		line := &order.Lines[index]
		if line.ReservedQuantity == 0 {
			continue
		}

		stock := store.Stocks[line.SKU]
		stock.Reserved -= line.ReservedQuantity
		store.Stocks[line.SKU] = stock
		line.ReservedQuantity = 0
	}

	order.Status = data.OrderStatusCancelled
	order.CancelledBy = cancelledBy
	order.CancellationReason = reason
	store.Orders[order.ID] = order

	return order, nil
}
