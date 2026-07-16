package orders

import (
	"fmt"

	"component-based-architecture/internal/components/quotes"
)

// Component owns order creation and order state for this lesson.
type Component struct {
	quotes quotes.ApprovedQuoteSource
	orders map[string]Order
	nextID int
}

func NewComponent(quotes quotes.ApprovedQuoteSource) *Component {
	return &Component{quotes: quotes, orders: make(map[string]Order)}
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

func (c *Component) ConvertQuoteToOrder(command ConvertQuoteToOrderCommand) (ConvertQuoteToOrderResult, error) {
	quote, err := c.quotes.GetApprovedQuoteForOrder(command.QuoteID)
	if err != nil {
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
