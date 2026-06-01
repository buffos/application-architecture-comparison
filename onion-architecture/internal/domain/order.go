package domain

import (
	"errors"
	"fmt"
	"sync/atomic"
)

var ErrOrderNotFound = errors.New("order not found")
var ErrOrderNotPayable = errors.New("order is not payable")
var ErrOrderNotShippable = errors.New("order is not shippable")

const OrderStatusPendingPayment = "PendingPayment"
const OrderStatusPaid = "Paid"
const OrderStatusShipped = "Shipped"

var orderSequence uint64

type OrderLine struct {
	ProductSKU      string
	ProductName     string
	ProductCategory string
	Quantity        int
	UnitPrice       int
}

type Order struct {
	ID         string
	QuoteID    string
	CustomerID string
	Status     string
	Lines      []OrderLine
}

func NewOrderFromQuote(quote Quote) (Order, error) {
	if err := quote.EnsureConvertible(); err != nil {
		return Order{}, err
	}

	id := atomic.AddUint64(&orderSequence, 1)
	lines := make([]OrderLine, 0, len(quote.Lines))
	for _, line := range quote.Lines {
		lines = append(lines, OrderLine{
			ProductSKU:      line.ProductSKU,
			ProductName:     line.ProductName,
			ProductCategory: line.ProductCategory,
			Quantity:        line.Quantity,
			UnitPrice:       line.UnitPrice,
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

func (o Order) EnsureShippable() error {
	if o.Status != OrderStatusPaid {
		return ErrOrderNotShippable
	}

	return nil
}

func (o *Order) MarkShipped() error {
	if err := o.EnsureShippable(); err != nil {
		return err
	}

	o.Status = OrderStatusShipped
	return nil
}
