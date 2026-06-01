package domain

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
var ErrOrderNotCancellable = errors.New("order is not cancellable")
var ErrOrderNotReturnable = errors.New("order is not returnable")
var ErrReturnLineNotFound = errors.New("return line not found")
var ErrReturnQuantityMustBePositive = errors.New("return quantity must be positive")
var ErrReturnQuantityExceedsReturnable = errors.New("return quantity exceeds returnable quantity")
var ErrReturnRequestHasNoLines = errors.New("return request has no lines")

const OrderStatusPendingPayment = "PendingPayment"
const OrderStatusPaymentReview = "PaymentReview"
const OrderStatusPaid = "Paid"
const OrderStatusPartiallyShipped = "PartiallyShipped"
const OrderStatusShipped = "Shipped"
const OrderStatusCancelled = "Cancelled"

var orderSequence uint64

type OrderLine struct {
	ProductSKU       string
	ProductName      string
	ProductCategory  string
	Quantity         int
	ShippedQuantity  int
	ReturnedQuantity int
	UnitPrice        int
	ReturnWindowDays int
}

type Order struct {
	ID         string
	QuoteID    string
	CustomerID string
	Status     string
	Lines      []OrderLine
	ShippedAt  time.Time
}

func NewOrderFromQuote(quote Quote) (Order, error) {
	if err := quote.EnsureConvertible(); err != nil {
		return Order{}, err
	}

	id := atomic.AddUint64(&orderSequence, 1)
	lines := make([]OrderLine, 0, len(quote.Lines))
	for _, line := range quote.Lines {
		lines = append(lines, OrderLine{
			ProductSKU:       line.ProductSKU,
			ProductName:      line.ProductName,
			ProductCategory:  line.ProductCategory,
			Quantity:         line.Quantity,
			UnitPrice:        line.UnitPrice,
			ReturnWindowDays: line.ReturnWindowDays,
		})
	}

	return Order{
		ID:         fmt.Sprintf("order-%03d", id),
		QuoteID:    quote.ID,
		CustomerID: quote.CustomerID,
		Status:     OrderStatusPendingPayment,
		Lines:      lines,
	}, nil
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

func (o Order) EnsureShippable() error {
	if o.Status != OrderStatusPaid && o.Status != OrderStatusPartiallyShipped {
		return ErrOrderNotShippable
	}

	return nil
}

func (o *Order) ApplyShipment(lines []ShipmentLine, shippedAt time.Time) error {
	resolved, err := resolveShipmentLines(*o, lines)
	if err != nil {
		return err
	}

	shipmentsBySKU := make(map[string]int, len(resolved))
	for _, line := range resolved {
		shipmentsBySKU[line.ProductSKU] += line.Quantity
	}

	for i := range o.Lines {
		o.Lines[i].ShippedQuantity += shipmentsBySKU[o.Lines[i].ProductSKU]
	}

	if o.ShippedAt.IsZero() {
		o.ShippedAt = shippedAt
	}

	allShipped := true
	for _, line := range o.Lines {
		if line.ShippedQuantity < line.Quantity {
			allShipped = false
			break
		}
	}

	if allShipped {
		o.Status = OrderStatusShipped
		return nil
	}

	o.Status = OrderStatusPartiallyShipped
	return nil
}

func (o *Order) MarkShipped(shippedAt time.Time) error {
	if err := o.EnsureShippable(); err != nil {
		return err
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

func (o Order) EnsureReturnable() error {
	if o.Status != OrderStatusShipped && o.Status != OrderStatusPartiallyShipped {
		return ErrOrderNotReturnable
	}

	for _, line := range o.Lines {
		if line.ShippedQuantity-line.ReturnedQuantity > 0 {
			return nil
		}
	}

	return ErrOrderNotReturnable
}

func (o Order) ResolveReturnLines(requested []ReturnRequestLine) ([]ReturnRequestLine, error) {
	if err := o.EnsureReturnable(); err != nil {
		return nil, err
	}

	if len(requested) == 0 {
		return nil, ErrReturnRequestHasNoLines
	}

	remainingBySKU := make(map[string]int, len(o.Lines))
	linesBySKU := make(map[string]OrderLine, len(o.Lines))
	for _, line := range o.Lines {
		remaining := line.ShippedQuantity - line.ReturnedQuantity
		if remaining > 0 {
			remainingBySKU[line.ProductSKU] = remaining
			linesBySKU[line.ProductSKU] = line
		}
	}

	accumulated := make(map[string]int)
	result := make([]ReturnRequestLine, 0, len(requested))
	for _, requestedLine := range requested {
		if requestedLine.Quantity <= 0 {
			return nil, ErrReturnQuantityMustBePositive
		}

		line, ok := linesBySKU[requestedLine.ProductSKU]
		if !ok {
			return nil, ErrReturnLineNotFound
		}

		accumulated[requestedLine.ProductSKU] += requestedLine.Quantity
		if accumulated[requestedLine.ProductSKU] > remainingBySKU[requestedLine.ProductSKU] {
			return nil, ErrReturnQuantityExceedsReturnable
		}

		result = append(result, ReturnRequestLine{
			ProductSKU:       line.ProductSKU,
			ProductCategory:  line.ProductCategory,
			Quantity:         requestedLine.Quantity,
			ReturnWindowDays: line.ReturnWindowDays,
		})
	}

	return result, nil
}

func (o *Order) ApplyReturn(lines []ReturnRequestLine) error {
	resolved, err := o.ResolveReturnLines(lines)
	if err != nil {
		return err
	}

	returnedBySKU := make(map[string]int, len(resolved))
	for _, line := range resolved {
		returnedBySKU[line.ProductSKU] += line.Quantity
	}

	for i := range o.Lines {
		o.Lines[i].ReturnedQuantity += returnedBySKU[o.Lines[i].ProductSKU]
	}

	return nil
}
