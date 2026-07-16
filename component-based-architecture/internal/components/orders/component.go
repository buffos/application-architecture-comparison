package orders

import (
	"fmt"

	"component-based-architecture/internal/components/inventory"
	"component-based-architecture/internal/components/payments"
	"component-based-architecture/internal/components/quotes"
)

// Component owns order creation and order state for this lesson.
type Component struct {
	quotes   quotes.ApprovedQuoteSource
	stock    inventory.Reserver
	payments payments.Processor
	orders   map[string]Order
	nextID   int
}

func NewComponent(quotes quotes.ApprovedQuoteSource, stock inventory.Reserver, payments payments.Processor) *Component {
	return &Component{quotes: quotes, stock: stock, payments: payments, orders: make(map[string]Order)}
}

type ConvertQuoteToOrderCommand struct {
	QuoteID string
}

type ConvertQuoteToOrderResult struct {
	OrderID    string
	QuoteID    string
	CustomerID string
	Status     string
	LineCount  int
}

type CapturePaymentCommand struct {
	OrderID string
}

type CapturePaymentResult struct {
	OrderID    string
	CustomerID string
	Status     string
	LineCount  int
}

func (c *Component) ConvertQuoteToOrder(command ConvertQuoteToOrderCommand) (ConvertQuoteToOrderResult, error) {
	quote, err := c.quotes.GetApprovedQuoteForOrder(command.QuoteID)
	if err != nil {
		return ConvertQuoteToOrderResult{}, err
	}

	reservations := make([]inventory.ReservationItem, 0, len(quote.Lines))
	for _, line := range quote.Lines {
		reservations = append(reservations, inventory.ReservationItem{ProductSKU: line.ProductSKU, Quantity: line.Quantity})
	}
	if err := c.stock.Reserve(reservations); err != nil {
		return ConvertQuoteToOrderResult{}, err
	}

	c.nextID++
	order := newOrderFromApprovedQuote(fmt.Sprintf("order-%03d", c.nextID), quote)
	c.orders[order.ID] = order

	return ConvertQuoteToOrderResult{
		OrderID: order.ID, QuoteID: order.QuoteID, CustomerID: order.CustomerID,
		Status: order.Status, LineCount: len(order.Lines),
	}, nil
}

func (c *Component) CapturePayment(command CapturePaymentCommand) (CapturePaymentResult, error) {
	order, ok := c.orders[command.OrderID]
	if !ok {
		return CapturePaymentResult{}, ErrOrderNotFound
	}

	amount := 0
	for _, line := range order.Lines {
		amount += line.Quantity * line.UnitPrice
	}
	if _, err := c.payments.Capture(payments.PaymentRequest{OrderID: order.ID, CustomerID: order.CustomerID, Amount: amount}); err != nil {
		return CapturePaymentResult{}, err
	}
	if err := order.MarkPaid(); err != nil {
		return CapturePaymentResult{}, err
	}
	c.orders[order.ID] = order

	return CapturePaymentResult{OrderID: order.ID, CustomerID: order.CustomerID, Status: order.Status, LineCount: len(order.Lines)}, nil
}
