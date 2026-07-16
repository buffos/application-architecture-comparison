package orders

import (
	"errors"
	"time"

	"component-based-architecture/internal/components/quotes"
)

var (
	ErrOrderNotFound       = errors.New("order not found")
	ErrOrderNotPayable     = errors.New("order is not payable")
	ErrOrderNotShippable   = errors.New("order is not shippable")
	ErrOrderNotCancellable = errors.New("order is not cancellable")
	ErrOrderNotReturnable  = errors.New("order is not returnable")
)

const (
	OrderStatusPendingPayment = "PendingPayment"
	OrderStatusPaid           = "Paid"
	OrderStatusShipped        = "Shipped"
	OrderStatusCancelled      = "Cancelled"
)

type Order struct {
	ID         string
	QuoteID    string
	CustomerID string
	Status     string
	Lines      []OrderLine
	ShippedAt  time.Time
}

type OrderLine struct {
	ProductSKU       string
	ProductName      string
	ProductCategory  string
	Quantity         int
	UnitPrice        int
	ReturnWindowDays int
}

func (o *Order) MarkPaid() error {
	if o.Status != OrderStatusPendingPayment {
		return ErrOrderNotPayable
	}
	o.Status = OrderStatusPaid
	return nil
}

func (o *Order) MarkShipped(shippedAt time.Time) error {
	if o.Status != OrderStatusPaid {
		return ErrOrderNotShippable
	}
	o.Status = OrderStatusShipped
	o.ShippedAt = shippedAt
	return nil
}

func (o *Order) Cancel() error {
	if o.Status == OrderStatusShipped || o.Status == OrderStatusCancelled {
		return ErrOrderNotCancellable
	}
	o.Status = OrderStatusCancelled
	return nil
}

func (o Order) EnsureReturnable() error {
	if o.Status != OrderStatusShipped {
		return ErrOrderNotReturnable
	}
	return nil
}

func newOrderFromApprovedQuote(id string, quote quotes.ApprovedQuote) Order {
	lines := make([]OrderLine, 0, len(quote.Lines))
	for _, line := range quote.Lines {
		lines = append(lines, OrderLine{
			ProductSKU: line.ProductSKU, ProductName: line.ProductName, ProductCategory: line.ProductCategory,
			Quantity: line.Quantity, UnitPrice: line.UnitPrice, ReturnWindowDays: line.ReturnWindowDays,
		})
	}
	return Order{
		ID: id, QuoteID: quote.QuoteID, CustomerID: quote.CustomerID, Status: OrderStatusPendingPayment, Lines: lines,
	}
}
