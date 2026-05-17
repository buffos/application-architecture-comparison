package domain

import (
	"errors"
	"fmt"
	"sync/atomic"
)

const OrderStatusReadyForPayment = "ReadyForPayment"
const OrderStatusReadyForFulfillment = "ReadyForFulfillment"
const OrderStatusShipped = "Shipped"

var orderSequence uint64

var ErrOrderNotFound = errors.New("order not found")
var ErrQuoteNotConvertible = errors.New("quote must be approved before conversion")
var ErrPaymentAlreadyAccepted = errors.New("payment is already accepted")

type OrderLine struct {
	SKU                 string
	ProductNameSnapshot string
	Quantity            int
}

type Order struct {
	ID            string
	SourceQuoteID string
	CustomerID    string
	Status        string
	PaymentStatus string
	Lines         []OrderLine
}

func NewOrderFromQuote(quote Quote) (Order, error) {
	if quote.Status != QuoteStatusApproved {
		return Order{}, ErrQuoteNotConvertible
	}

	id := atomic.AddUint64(&orderSequence, 1)
	lines := make([]OrderLine, 0, len(quote.Lines))

	for _, line := range quote.Lines {
		lines = append(lines, OrderLine{
			SKU:                 line.SKU,
			ProductNameSnapshot: line.ProductNameSnapshot,
			Quantity:            line.Quantity,
		})
	}

	return Order{
		ID:            fmt.Sprintf("order-%03d", id),
		SourceQuoteID: quote.ID,
		CustomerID:    quote.CustomerID,
		Status:        OrderStatusReadyForPayment,
		PaymentStatus: "Pending",
		Lines:         lines,
	}, nil
}

func (o *Order) AcceptPayment() error {
	if o.PaymentStatus == "Accepted" {
		return ErrPaymentAlreadyAccepted
	}

	o.PaymentStatus = "Accepted"
	o.Status = OrderStatusReadyForFulfillment
	return nil
}

func (o *Order) MarkShipped() error {
	if o.Status != OrderStatusReadyForFulfillment {
		return ErrShipmentNotAllowedUntilPaymentAccepted
	}

	o.Status = OrderStatusShipped
	return nil
}
