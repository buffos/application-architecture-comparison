package entities

import (
	"fmt"
	"sync/atomic"
	"time"
)

const OrderStatusPendingPayment = "PendingPayment"
const OrderStatusPaymentReview = "PaymentReview"
const OrderStatusPaid = "Paid"
const OrderStatusPartiallyShipped = "PartiallyShipped"
const OrderStatusShipped = "Shipped"
const OrderStatusCancelled = "Cancelled"

var orderSequence uint64

var ErrQuoteNotConvertible = ErrQuoteCannotTransition

type OrderLine struct {
	SKU              string
	ProductName      string
	Quantity         int
	ShippedQuantity  int
	ReturnedQuantity int
	UnitPrice        int
	LineTotal        int
	ReturnWindowDays int
}

type Order struct {
	ID            string
	CustomerID    string
	SourceQuoteID string
	Status        string
	Lines         []OrderLine
	ShippedAt     *time.Time
}

func NewOrderFromApprovedQuote(quote Quote, now time.Time) (Order, error) {
	if quote.Status != QuoteStatusApproved {
		return Order{}, ErrQuoteNotConvertible
	}

	id := atomic.AddUint64(&orderSequence, 1)
	lines := make([]OrderLine, 0, len(quote.Lines))
	for _, line := range quote.Lines {
		lines = append(lines, OrderLine{
			SKU:              line.SKU,
			ProductName:      line.ProductName,
			Quantity:         line.Quantity,
			UnitPrice:        line.UnitPrice,
			LineTotal:        line.LineTotal,
			ReturnWindowDays: line.ReturnWindowDays,
		})
	}

	return Order{
		ID:            fmt.Sprintf("order-%03d", id),
		CustomerID:    quote.CustomerID,
		SourceQuoteID: quote.ID,
		Status:        OrderStatusPendingPayment,
		Lines:         lines,
	}, nil
}

func (o *Order) MarkPaid() error {
	if o.Status != OrderStatusPendingPayment && o.Status != OrderStatusPaymentReview {
		return ErrQuoteCannotTransition
	}

	o.Status = OrderStatusPaid
	return nil
}

func (o *Order) MarkPaymentReview() error {
	if o.Status != OrderStatusPendingPayment {
		return ErrQuoteCannotTransition
	}

	o.Status = OrderStatusPaymentReview
	return nil
}

func (o *Order) ApprovePaymentReview() error {
	if o.Status != OrderStatusPaymentReview {
		return ErrQuoteCannotTransition
	}

	o.Status = OrderStatusPaid
	return nil
}

func (o *Order) MarkShipped() error {
	return o.MarkShippedAt(time.Now())
}

func (o *Order) MarkShippedAt(at time.Time) error {
	if o.Status != OrderStatusPaid {
		return ErrQuoteCannotTransition
	}

	lines := make([]ShipmentLine, 0, len(o.Lines))
	for _, line := range o.Lines {
		lines = append(lines, ShipmentLine{
			SKU:         line.SKU,
			ProductName: line.ProductName,
			Quantity:    line.Quantity - line.ShippedQuantity,
		})
	}

	return o.ApplyShipment(lines, at)
}

func (o *Order) ApplyShipment(lines []ShipmentLine, at time.Time) error {
	if o.Status != OrderStatusPaid && o.Status != OrderStatusPartiallyShipped {
		return ErrQuoteCannotTransition
	}

	if len(lines) == 0 {
		return ErrQuoteCannotTransition
	}

	for _, shipmentLine := range lines {
		if shipmentLine.Quantity <= 0 {
			return ErrQuoteCannotTransition
		}

		matched := false
		for idx := range o.Lines {
			orderLine := &o.Lines[idx]
			if orderLine.SKU != shipmentLine.SKU {
				continue
			}

			remaining := orderLine.Quantity - orderLine.ShippedQuantity
			if shipmentLine.Quantity > remaining {
				return ErrQuoteCannotTransition
			}

			matched = true
			break
		}

		if !matched {
			return ErrQuoteCannotTransition
		}
	}

	for _, shipmentLine := range lines {
		for idx := range o.Lines {
			orderLine := &o.Lines[idx]
			if orderLine.SKU == shipmentLine.SKU {
				orderLine.ShippedQuantity += shipmentLine.Quantity
				break
			}
		}
	}

	allShipped := true
	for _, line := range o.Lines {
		if line.ShippedQuantity < line.Quantity {
			allShipped = false
			break
		}
	}

	copy := at
	o.ShippedAt = &copy
	if allShipped {
		o.Status = OrderStatusShipped
	} else {
		o.Status = OrderStatusPartiallyShipped
	}

	return nil
}

func (o *Order) Cancel() error {
	if o.Status == OrderStatusShipped || o.Status == OrderStatusPartiallyShipped {
		return ErrQuoteCannotTransition
	}

	if o.Status == OrderStatusCancelled {
		return ErrQuoteCannotTransition
	}

	o.Status = OrderStatusCancelled
	return nil
}

func (o *Order) ApplyReturn(lines []ReturnRequestLine) error {
	if o.Status != OrderStatusShipped && o.Status != OrderStatusPartiallyShipped {
		return ErrQuoteCannotTransition
	}

	if len(lines) == 0 {
		return ErrQuoteCannotTransition
	}

	for _, returnLine := range lines {
		if returnLine.Quantity <= 0 {
			return ErrQuoteCannotTransition
		}

		matched := false
		for idx := range o.Lines {
			orderLine := &o.Lines[idx]
			if orderLine.SKU != returnLine.SKU {
				continue
			}

			returnable := orderLine.ShippedQuantity - orderLine.ReturnedQuantity
			if returnLine.Quantity > returnable {
				return ErrQuoteCannotTransition
			}

			matched = true
			break
		}

		if !matched {
			return ErrQuoteCannotTransition
		}
	}

	for _, returnLine := range lines {
		for idx := range o.Lines {
			orderLine := &o.Lines[idx]
			if orderLine.SKU == returnLine.SKU {
				orderLine.ReturnedQuantity += returnLine.Quantity
				break
			}
		}
	}

	return nil
}
