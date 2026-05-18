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
var ErrQuoteNotApproved = errors.New("quote must be approved before conversion")
var ErrPaymentFailed = errors.New("payment failed")

type OrderLine struct {
	SKU               string
	ProductName       string
	ProductCategory   string
	Quantity          int
	BaseUnitPrice     int
	AdjustedUnitPrice int
	LineTotal         int
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
		return Order{}, ErrQuoteNotApproved
	}

	id := atomic.AddUint64(&orderSequence, 1)
	lines := make([]OrderLine, 0, len(quote.Lines))

	for _, line := range quote.Lines {
		lines = append(lines, OrderLine{
			SKU:               line.SKU,
			ProductName:       line.ProductName,
			ProductCategory:   line.ProductCategory,
			Quantity:          line.Quantity,
			BaseUnitPrice:     line.BaseUnitPrice,
			AdjustedUnitPrice: line.AdjustedUnitPrice,
			LineTotal:         line.LineTotal,
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

func (o *Order) AcceptPayment() {
	o.PaymentStatus = "Accepted"
	o.Status = OrderStatusReadyForFulfillment
}

func (o *Order) MarkShipped() error {
	if o.Status != OrderStatusReadyForFulfillment {
		return ErrShipmentNotAllowedUntilPaymentAccepted
	}

	o.Status = OrderStatusShipped
	return nil
}
