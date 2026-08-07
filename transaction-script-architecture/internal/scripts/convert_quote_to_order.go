package scripts

import (
	"errors"
	"fmt"

	"transaction-script-architecture/internal/data"
)

var (
	ErrQuoteNotConvertible   = errors.New("quote must be approved")
	ErrQuoteAlreadyConverted = errors.New("quote is already converted")
	ErrRequestedByRequired   = errors.New("requester is required")
)

// ConvertQuoteToOrder copies an approved quote into a new passive order. It
// deliberately stops before inventory reservation; that coordination arrives
// in the next lesson.
func ConvertQuoteToOrder(store *data.Store, quoteID string, requestedBy string) (data.Order, error) {
	if store == nil {
		return data.Order{}, ErrStoreRequired
	}

	if quoteID == "" {
		return data.Order{}, ErrQuoteIDRequired
	}

	if requestedBy == "" {
		return data.Order{}, ErrRequestedByRequired
	}

	quote, ok := store.Quotes[quoteID]
	if !ok {
		return data.Order{}, ErrQuoteNotFound
	}

	if quote.Status == data.QuoteStatusConverted {
		return data.Order{}, ErrQuoteAlreadyConverted
	}

	if quote.Status != data.QuoteStatusApproved {
		return data.Order{}, ErrQuoteNotConvertible
	}

	store.NextOrderNumber++
	order := data.Order{
		ID:            fmt.Sprintf("order-%03d", store.NextOrderNumber),
		SourceQuoteID: quote.ID,
		CustomerID:    quote.CustomerID,
		Status:        data.OrderStatusPendingReservation,
		RequestedBy:   requestedBy,
		PaymentStatus: data.PaymentStatusNotRequired,
		Lines:         make([]data.OrderLine, 0, len(quote.Lines)),
	}

	for index, line := range quote.Lines {
		order.Lines = append(order.Lines, data.OrderLine{
			ID:                  fmt.Sprintf("%s-line-%03d", order.ID, index+1),
			SKU:                 line.SKU,
			ProductNameSnapshot: line.ProductNameSnapshot,
			ProductCategory:     line.ProductCategory,
			OrderedQuantity:     line.Quantity,
			UnitPrice:           line.UnitPrice,
			DiscountAmount:      line.DiscountAmount,
			LineTotal:           line.LineTotal,
		})
		order.Total += line.LineTotal
	}

	store.Orders[order.ID] = order
	quote.Status = data.QuoteStatusConverted
	quote.ConvertedOrderID = order.ID
	store.Quotes[quote.ID] = quote

	return order, nil
}
