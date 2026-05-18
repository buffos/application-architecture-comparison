package domain

import (
	"errors"
	"fmt"
	"sync/atomic"
)

const OrderStatusReadyForPayment = "ReadyForPayment"
const OrderStatusReadyForFulfillment = "ReadyForFulfillment"
const OrderStatusShipped = "Shipped"
const OrderStatusCancelled = "Cancelled"

var orderSequence uint64

var ErrOrderNotFound = errors.New("order not found")
var ErrQuoteNotConvertible = errors.New("quote must be approved before conversion")
var ErrPaymentAlreadyAccepted = errors.New("payment is already accepted")
var ErrOrderAlreadyCancelled = errors.New("order is already cancelled")
var ErrOrderAlreadyShipped = errors.New("order is already shipped")

type OrderLine struct {
	SKU                 string
	ProductCategory     string
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
			ProductCategory:     line.ProductCategory,
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

func (o *Order) Cancel() error {
	if o.Status == OrderStatusCancelled {
		return ErrOrderAlreadyCancelled
	}

	if o.Status == OrderStatusShipped {
		return ErrOrderAlreadyShipped
	}

	o.Status = OrderStatusCancelled
	return nil
}
