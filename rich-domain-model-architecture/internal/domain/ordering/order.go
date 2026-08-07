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
	OrderStatusPendingPayment   OrderStatus = "PendingPayment"
	OrderStatusPaymentReview    OrderStatus = "PaymentReview"
	OrderStatusPartiallyShipped OrderStatus = "PartiallyShipped"
	OrderStatusPaid             OrderStatus = "Paid"
	OrderStatusShipped          OrderStatus = "Shipped"
	OrderStatusCancelled        OrderStatus = "Cancelled"
)

type OrderLine struct {
	sku             ProductSKU
	productName     string
	category        ProductCategory
	quantity        int
	unitPrice       Money
	shippedQuantity int
}

func (line OrderLine) ProductSKU() ProductSKU           { return line.sku }
func (line OrderLine) ProductName() string              { return line.productName }
func (line OrderLine) ProductCategory() ProductCategory { return line.category }
func (line OrderLine) Quantity() int                    { return line.quantity }
func (line OrderLine) UnitPrice() Money                 { return line.unitPrice }
func (line OrderLine) ShippedQuantity() int             { return line.shippedQuantity }

type ShipmentSelection struct {
	ProductSKU ProductSKU
	Quantity   int
}

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

func (order Order) ID() OrderID            { return order.id }
func (order Order) QuoteID() QuoteID       { return order.quoteID }
func (order Order) CustomerID() CustomerID { return order.customerID }
func (order Order) Status() OrderStatus    { return order.status }
func (order Order) Lines() []OrderLine     { return append([]OrderLine(nil), order.lines...) }

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

func (order *Order) MarkPaymentReview() error {
	if order.status != OrderStatusPendingPayment {
		return ErrOrderNotReviewable
	}
	order.status = OrderStatusPaymentReview
	return nil
}

func (order *Order) ApprovePaymentReview() error {
	if order.status != OrderStatusPaymentReview {
		return ErrOrderNotReviewable
	}
	order.status = OrderStatusPaid
	return nil
}

func (order *Order) MarkShipped() error {
	return order.ApplyShipment(nil)
}

func (order *Order) ApplyShipment(selections []ShipmentSelection) error {
	if order.status != OrderStatusPaid && order.status != OrderStatusPartiallyShipped {
		return ErrOrderNotShippable
	}
	if len(selections) == 0 {
		selections = make([]ShipmentSelection, 0, len(order.lines))
		for _, line := range order.lines {
			if remaining := line.quantity - line.shippedQuantity; remaining > 0 {
				selections = append(selections, ShipmentSelection{ProductSKU: line.sku, Quantity: remaining})
			}
		}
	}

	requested := make(map[ProductSKU]int, len(selections))
	for _, selection := range selections {
		if selection.Quantity <= 0 {
			return ErrOrderNotShippable
		}
		requested[selection.ProductSKU] += selection.Quantity
	}
	for sku, quantity := range requested {
		matched := false
		for _, line := range order.lines {
			if line.sku != sku {
				continue
			}
			if quantity > line.quantity-line.shippedQuantity {
				return ErrOrderNotShippable
			}
			matched = true
			break
		}
		if !matched {
			return ErrOrderNotShippable
		}
	}

	for sku, quantity := range requested {
		for index := range order.lines {
			if order.lines[index].sku == sku {
				order.lines[index].shippedQuantity += quantity
			}
		}
	}
	for _, line := range order.lines {
		if line.shippedQuantity < line.quantity {
			order.status = OrderStatusPartiallyShipped
			return nil
		}
	}
	order.status = OrderStatusShipped
	return nil
}

func (order *Order) Cancel() error {
	if order.status != OrderStatusPendingPayment && order.status != OrderStatusPaid {
		return ErrOrderNotCancellable
	}
	order.status = OrderStatusCancelled
	return nil
}
