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
	ErrInsufficientStock     = errors.New("insufficient stock")
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

	if err := reserveOrderStock(store, &order); err != nil {
		return data.Order{}, err
	}

	store.Orders[order.ID] = order
	quote.Status = data.QuoteStatusConverted
	quote.ConvertedOrderID = order.ID
	store.Quotes[quote.ID] = quote

	return order, nil
}

type stockReservation struct {
	sku      string
	quantity int
}

func reserveOrderStock(store *data.Store, order *data.Order) error {
	reservations := make([]stockReservation, 0, len(order.Lines))
	backordered := false

	for _, line := range order.Lines {
		stock, ok := store.Stocks[line.SKU]
		available := 0
		if ok {
			available = stock.OnHand - stock.Reserved
		}

		if !ok || available < line.OrderedQuantity {
			policy := data.StockShortageRejectOrder
			if product, exists := store.Products[line.SKU]; exists && product.StockShortagePolicy != "" {
				policy = product.StockShortagePolicy
			}

			if policy == data.StockShortageAllowBackorder {
				backordered = true
				continue
			}

			return ErrInsufficientStock
		}

		reservations = append(reservations, stockReservation{
			sku:      line.SKU,
			quantity: line.OrderedQuantity,
		})
	}

	if store.Stocks == nil {
		store.Stocks = make(map[string]data.StockRecord)
	}

	for _, reservation := range reservations {
		stock := store.Stocks[reservation.sku]
		stock.Reserved += reservation.quantity
		store.Stocks[reservation.sku] = stock
		for index := range order.Lines {
			if order.Lines[index].SKU == reservation.sku {
				order.Lines[index].ReservedQuantity = reservation.quantity
				break
			}
		}
	}

	if backordered {
		order.Status = data.OrderStatusBackordered
	} else {
		order.Status = data.OrderStatusReadyForPayment
		order.PaymentStatus = data.PaymentStatusPending
	}

	return nil
}
