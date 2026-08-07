package ordering

import (
	"errors"

	"rich-domain-model-architecture/internal/domain/quoting"
)

var (
	ErrOrderIDRequired         = errors.New("order id is required")
	ErrQuoteNotApproved        = errors.New("quote is not approved")
	ErrOrderHasNoLines         = errors.New("order must contain at least one line")
	ErrOrderNotAwaitingPayment = errors.New("order is not awaiting payment")
)

type OrderID string
type QuoteID string
type CustomerID string
type ProductSKU string
type ProductCategory string

type OrderStatus string

const (
	OrderStatusPendingPayment OrderStatus = "PendingPayment"
	OrderStatusPaid           OrderStatus = "Paid"
)

type OrderLine struct {
	sku         ProductSKU
	productName string
	category    ProductCategory
	quantity    int
	unitPrice   Money
}

func (line OrderLine) ProductSKU() ProductSKU         { return line.sku }
func (line OrderLine) ProductName() string             { return line.productName }
func (line OrderLine) ProductCategory() ProductCategory { return line.category }
func (line OrderLine) Quantity() int                   { return line.quantity }
func (line OrderLine) UnitPrice() Money                { return line.unitPrice }

// Order is the aggregate root for committed commercial transactions.
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
		lines = append(lines, OrderLine{
			sku:         ProductSKU(line.ProductSKU()),
			productName: line.ProductName(),
			category:    ProductCategory(line.ProductCategory()),
			quantity:    line.Quantity(),
			unitPrice:   price,
		})
	}

	return Order{
		id:         id,
		quoteID:    QuoteID(quote.ID()),
		customerID: CustomerID(quote.CustomerID()),
		status:     OrderStatusPendingPayment,
		lines:      lines,
	}, nil
}

func (order Order) ID() OrderID             { return order.id }
func (order Order) QuoteID() QuoteID        { return order.quoteID }
func (order Order) CustomerID() CustomerID  { return order.customerID }
func (order Order) Status() OrderStatus     { return order.status }
func (order Order) Lines() []OrderLine      { return append([]OrderLine(nil), order.lines...) }

func (order Order) Total() (Money, error) {
	if len(order.lines) == 0 {
		return Money{}, ErrOrderHasNoLines
	}
	total, err := NewMoney(0, order.lines[0].unitPrice.Currency())
	if err != nil {
		return Money{}, err
	}
	for _, line := range order.lines {
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

func (order *Order) MarkPaid() error {
	if order.status != OrderStatusPendingPayment {
		return ErrOrderNotAwaitingPayment
	}
	order.status = OrderStatusPaid
	return nil
}
