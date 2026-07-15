package orders

import (
	"errors"
	"fmt"
	"sync/atomic"
	"time"
)

var ErrOrderNotFound = errors.New("order not found")
var ErrOrderNotPayable = errors.New("order is not payable")
var ErrOrderNotPaymentReviewable = errors.New("order is not awaiting payment review")
var ErrOrderNotShippable = errors.New("order is not shippable")
var ErrShipmentQuantityInvalid = errors.New("shipment quantity is invalid")
var ErrOrderNotCancellable = errors.New("order is not cancellable")
var ErrOrderNotReturnable = errors.New("order is not returnable")

const OrderStatusPendingPayment = "PendingPayment"
const OrderStatusPaymentReview = "PaymentReview"
const OrderStatusPaid = "Paid"
const OrderStatusPartiallyShipped = "PartiallyShipped"
const OrderStatusShipped = "Shipped"
const OrderStatusCancelled = "Cancelled"

var orderSequence uint64

type Order struct {
	ID         string
	QuoteID    string
	CustomerID string
	Status     string
	ShippedAt  time.Time
	Lines      []OrderLine
}

type OrderLine struct {
	ProductSKU       string
	ProductName      string
	ProductCategory  string
	Quantity         int
	ShippedQuantity  int
	UnitPrice        int
	ReturnWindowDays int
}

type ShipmentSelection struct {
	ProductSKU string
	Quantity   int
}

func NewOrderFromApprovedQuote(quoteID string, customerID string, lines []OrderLine) Order {
	id := atomic.AddUint64(&orderSequence, 1)

	return Order{
		ID:         fmt.Sprintf("order-%03d", id),
		QuoteID:    quoteID,
		CustomerID: customerID,
		Status:     OrderStatusPendingPayment,
		Lines:      lines,
	}
}

func (o Order) TotalAmount() int {
	total := 0
	for _, line := range o.Lines {
		total += line.Quantity * line.UnitPrice
	}

	return total
}

func (o *Order) MarkPaid() error {
	if o.Status != OrderStatusPendingPayment {
		return ErrOrderNotPayable
	}

	o.Status = OrderStatusPaid
	return nil
}

func (o *Order) MarkPaymentReview() error {
	if o.Status != OrderStatusPendingPayment {
		return ErrOrderNotPayable
	}

	o.Status = OrderStatusPaymentReview
	return nil
}

func (o *Order) ApprovePaymentReview() error {
	if o.Status != OrderStatusPaymentReview {
		return ErrOrderNotPaymentReviewable
	}

	o.Status = OrderStatusPaid
	return nil
}

func (o *Order) MarkShipped(shippedAt time.Time) error {
	if o.Status != OrderStatusPaid && o.Status != OrderStatusPartiallyShipped {
		return ErrOrderNotShippable
	}

	o.Status = OrderStatusShipped
	o.ShippedAt = shippedAt
	return nil
}

func (o *Order) Cancel() error {
	if o.Status == OrderStatusPartiallyShipped || o.Status == OrderStatusShipped || o.Status == OrderStatusCancelled {
		return ErrOrderNotCancellable
	}

	o.Status = OrderStatusCancelled
	return nil
}

func (o Order) RemainingShipmentSelections() []ShipmentSelection {
	lines := make([]ShipmentSelection, 0, len(o.Lines))
	for _, line := range o.Lines {
		remaining := line.Quantity - line.ShippedQuantity
		if remaining > 0 {
			lines = append(lines, ShipmentSelection{
				ProductSKU: line.ProductSKU,
				Quantity:   remaining,
			})
		}
	}

	return lines
}

func (o *Order) ApplyShipment(selections []ShipmentSelection, shippedAt time.Time) error {
	if o.Status != OrderStatusPaid && o.Status != OrderStatusPartiallyShipped {
		return ErrOrderNotShippable
	}

	if len(selections) == 0 {
		return ErrShipmentQuantityInvalid
	}

	for _, selection := range selections {
		if selection.Quantity <= 0 {
			return ErrShipmentQuantityInvalid
		}

		matched := false
		for i := range o.Lines {
			line := &o.Lines[i]
			if line.ProductSKU != selection.ProductSKU {
				continue
			}

			remaining := line.Quantity - line.ShippedQuantity
			if selection.Quantity > remaining {
				return ErrShipmentQuantityInvalid
			}

			line.ShippedQuantity += selection.Quantity
			matched = true
			break
		}

		if !matched {
			return ErrShipmentQuantityInvalid
		}
	}

	if len(o.RemainingShipmentSelections()) == 0 {
		o.Status = OrderStatusShipped
		o.ShippedAt = shippedAt
		return nil
	}

	o.Status = OrderStatusPartiallyShipped
	return nil
}

func (o Order) EnsureReturnable() error {
	if o.Status != OrderStatusShipped {
		return ErrOrderNotReturnable
	}

	return nil
}
