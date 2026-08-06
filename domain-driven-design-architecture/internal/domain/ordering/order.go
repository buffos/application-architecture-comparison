package ordering

import (
	"errors"

	"domain-driven-design-architecture/internal/domain/quoting"
)

var (
	ErrOrderIDRequired         = errors.New("order id is required")
	ErrQuoteNotApproved        = errors.New("quote is not approved")
	ErrOrderHasNoLines         = errors.New("order must contain at least one line")
	ErrOrderNotAwaitingPayment = errors.New("order is not awaiting payment")
	ErrOrderNotShippable       = errors.New("order is not shippable")
	ErrOrderNotCancellable     = errors.New("order is not cancellable")
	ErrOrderNotReviewable      = errors.New("order is not in payment review")
)

type OrderID string
type QuoteID string
type CustomerID string
type ProductSKU string
type ProductCategory string

type OrderStatus string

const (
	OrderStatusPendingPayment OrderStatus = "PendingPayment"
	OrderStatusPaymentReview  OrderStatus = "PaymentReview"
	OrderStatusPaid           OrderStatus = "Paid"
	OrderStatusShipped        OrderStatus = "Shipped"
	OrderStatusCancelled      OrderStatus = "Cancelled"
)

type OrderLine struct {
	sku       ProductSKU
	category  ProductCategory
	quantity  int
	unitPrice Money
}

func (line OrderLine) ProductSKU() ProductSKU           { return line.sku }
func (line OrderLine) ProductCategory() ProductCategory { return line.category }
func (line OrderLine) Quantity() int                    { return line.quantity }
func (line OrderLine) UnitPrice() Money                 { return line.unitPrice }

// Order is the aggregate root for the Ordering bounded context.
type Order struct {
	id         OrderID
	quoteID    QuoteID
	customerID CustomerID
	status     OrderStatus
	lines      []OrderLine
}

func NewOrderFromQuote(id OrderID, quote quoting.Quote) (Order, error) {
	if id == "" {
		return Order{}, ErrOrderIDRequired
	}
	if quote.Status() != quoting.QuoteStatusApproved {
		return Order{}, ErrQuoteNotApproved
	}
	quoteLines := quote.Lines()
	if len(quoteLines) == 0 {
		return Order{}, ErrOrderHasNoLines
	}
	lines := make([]OrderLine, 0, len(quoteLines))
	for _, line := range quoteLines {
		price, err := NewMoney(line.UnitPrice().Cents(), line.UnitPrice().Currency())
		if err != nil {
			return Order{}, err
		}
		lines = append(lines, OrderLine{sku: ProductSKU(line.ProductSKU()), category: ProductCategory(line.ProductCategory()), quantity: line.Quantity(), unitPrice: price})
	}
	return Order{id: id, quoteID: QuoteID(quote.ID()), customerID: CustomerID(quote.CustomerID()), status: OrderStatusPendingPayment, lines: lines}, nil
}

func (o Order) ID() OrderID            { return o.id }
func (o Order) QuoteID() QuoteID       { return o.quoteID }
func (o Order) CustomerID() CustomerID { return o.customerID }
func (o Order) Status() OrderStatus    { return o.status }
func (o Order) Lines() []OrderLine     { return append([]OrderLine(nil), o.lines...) }

func (o Order) Total() (Money, error) {
	if len(o.lines) == 0 {
		return Money{}, ErrOrderHasNoLines
	}
	total, err := NewMoney(0, o.lines[0].unitPrice.Currency())
	if err != nil {
		return Money{}, err
	}
	for _, line := range o.lines {
		lineTotal, err := line.unitPrice.Multiply(line.quantity)
		if err != nil {
			return Money{}, err
		}
		total, err = total.Add(lineTotal)
		if err != nil {
			return Money{}, err
		}
	}
	return total, nil
}

func (o *Order) MarkPaid() error {
	if o.status != OrderStatusPendingPayment {
		return ErrOrderNotAwaitingPayment
	}
	o.status = OrderStatusPaid
	return nil
}

func (o *Order) MarkPaymentReview() error {
	if o.status != OrderStatusPendingPayment {
		return ErrOrderNotReviewable
	}
	o.status = OrderStatusPaymentReview
	return nil
}

func (o *Order) ApprovePaymentReview() error {
	if o.status != OrderStatusPaymentReview {
		return ErrOrderNotReviewable
	}
	o.status = OrderStatusPaid
	return nil
}

func (o *Order) MarkShipped() error {
	if o.status != OrderStatusPaid {
		return ErrOrderNotShippable
	}
	o.status = OrderStatusShipped
	return nil
}

func (o *Order) Cancel() error {
	if o.status != OrderStatusPendingPayment && o.status != OrderStatusPaid {
		return ErrOrderNotCancellable
	}
	o.status = OrderStatusCancelled
	return nil
}
