package domain

import (
	"errors"
	"fmt"
	"sync/atomic"
	"time"
)

const OrderStatusReadyForPayment = "ReadyForPayment"
const OrderStatusPaymentReview = "PaymentReview"
const OrderStatusReadyForFulfillment = "ReadyForFulfillment"
const OrderStatusPartiallyShipped = "PartiallyShipped"
const OrderStatusShipped = "Shipped"
const OrderStatusCancelled = "Cancelled"

var orderSequence uint64

var ErrOrderNotFound = errors.New("order not found")
var ErrQuoteNotApproved = errors.New("quote must be approved before conversion")
var ErrPaymentFailed = errors.New("payment failed")
var ErrPaymentReviewNotAllowed = errors.New("payment review is not allowed")
var ErrOrderActorRequired = errors.New("actor is required")
var ErrOrderCancellationNotAllowed = errors.New("order cancellation is not allowed after shipment")
var ErrShipmentLineInvalid = errors.New("shipment line is invalid")
var ErrShipmentQuantityExceedsRemaining = errors.New("shipment quantity exceeds remaining shippable quantity")

type OrderLine struct {
	SKU               string
	ProductName       string
	ProductCategory   string
	Quantity          int
	BaseUnitPrice     int
	AdjustedUnitPrice int
	LineTotal         int
	ReturnWindowDays  int
	ShippedQuantity   int
	ReturnedQuantity  int
}

type Order struct {
	ID                string
	SourceQuoteID     string
	CustomerID        string
	Status            string
	PaymentStatus     string
	PaymentReviewedBy string
	ShippedAt         time.Time
	Lines             []OrderLine
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
			ReturnWindowDays:  line.ReturnWindowDays,
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

func (o *Order) MarkPaymentReview() {
	o.PaymentStatus = "ManualReview"
	o.Status = OrderStatusPaymentReview
}

func (o *Order) FailPayment() {
	o.PaymentStatus = "Failed"
}

func (o *Order) ApprovePaymentReview(reviewedBy string) error {
	if o.Status != OrderStatusPaymentReview {
		return ErrPaymentReviewNotAllowed
	}
	if reviewedBy == "" {
		return ErrOrderActorRequired
	}

	o.PaymentReviewedBy = reviewedBy
	o.AcceptPayment()
	return nil
}

func (o *Order) MarkShipped(shippedAt time.Time) error {
	if o.Status != OrderStatusReadyForFulfillment {
		return ErrShipmentNotAllowedUntilPaymentAccepted
	}

	o.Status = OrderStatusShipped
	o.ShippedAt = shippedAt
	return nil
}

func (o *OrderLine) RemainingShippableQuantity() int {
	return o.Quantity - o.ShippedQuantity
}

func (o *OrderLine) RemainingReturnableQuantity() int {
	return o.ShippedQuantity - o.ReturnedQuantity
}

func (o *Order) ApplyShipment(lines []ShipmentLine, shippedAt time.Time) error {
	if o.Status != OrderStatusReadyForFulfillment && o.Status != OrderStatusPartiallyShipped {
		return ErrShipmentNotAllowedUntilPaymentAccepted
	}
	if len(lines) == 0 {
		return ErrShipmentLineInvalid
	}

	indexBySKU := make(map[string]int, len(o.Lines))
	for i, line := range o.Lines {
		indexBySKU[line.SKU] = i
	}

	for _, shippedLine := range lines {
		if shippedLine.Quantity <= 0 {
			return ErrShipmentLineInvalid
		}
		index, ok := indexBySKU[shippedLine.SKU]
		if !ok {
			return ErrShipmentLineInvalid
		}
		if shippedLine.Quantity > o.Lines[index].RemainingShippableQuantity() {
			return ErrShipmentQuantityExceedsRemaining
		}
	}

	for _, shippedLine := range lines {
		index := indexBySKU[shippedLine.SKU]
		o.Lines[index].ShippedQuantity += shippedLine.Quantity
	}

	fullyShipped := true
	for _, line := range o.Lines {
		if line.RemainingShippableQuantity() > 0 {
			fullyShipped = false
			break
		}
	}

	if fullyShipped {
		o.Status = OrderStatusShipped
		o.ShippedAt = shippedAt
		return nil
	}

	o.Status = OrderStatusPartiallyShipped
	if o.ShippedAt.IsZero() {
		o.ShippedAt = shippedAt
	}
	return nil
}

func (o *Order) ApplyAcceptedReturn(lines []ReturnLine) error {
	indexBySKU := make(map[string]int, len(o.Lines))
	for i, line := range o.Lines {
		indexBySKU[line.SKU] = i
	}

	for _, returnLine := range lines {
		index, ok := indexBySKU[returnLine.SKU]
		if !ok {
			return ErrReturnLineInvalid
		}
		if returnLine.Quantity <= 0 {
			return ErrReturnLineInvalid
		}
		if returnLine.Quantity > o.Lines[index].RemainingReturnableQuantity() {
			return ErrReturnQuantityExceedsRemaining
		}
	}

	for _, returnLine := range lines {
		index := indexBySKU[returnLine.SKU]
		o.Lines[index].ReturnedQuantity += returnLine.Quantity
	}

	return nil
}

func (o *Order) Cancel() error {
	if o.Status == OrderStatusShipped {
		return ErrOrderCancellationNotAllowed
	}

	o.Status = OrderStatusCancelled
	return nil
}
