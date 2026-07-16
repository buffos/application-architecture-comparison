package orders

import "time"

// Reader is the public read contract provided by Orders. It exposes business
// views rather than the component's private order map.
type Reader interface {
	GetOrder(query GetOrderQuery) (OrderDetails, error)
	ListOrders(query ListOrdersQuery) []OrderSummary
}

type GetOrderQuery struct {
	OrderID string
}

type ListOrdersQuery struct {
	Status string
}

type OrderDetails struct {
	OrderID    string
	QuoteID    string
	CustomerID string
	Status     string
	LineCount  int
	ShippedAt  time.Time
	Lines      []OrderLineDetails
}

type OrderLineDetails struct {
	ProductSKU  string
	ProductName string
	Quantity    int
	UnitPrice   int
}

type OrderSummary struct {
	OrderID    string
	QuoteID    string
	CustomerID string
	Status     string
	LineCount  int
}

func (c *Component) GetOrder(query GetOrderQuery) (OrderDetails, error) {
	order, ok := c.orders[query.OrderID]
	if !ok {
		return OrderDetails{}, ErrOrderNotFound
	}
	return orderDetails(order), nil
}

func (c *Component) ListOrders(query ListOrdersQuery) []OrderSummary {
	orders := make([]OrderSummary, 0, len(c.orders))
	for _, order := range c.orders {
		if query.Status != "" && order.Status != query.Status {
			continue
		}
		orders = append(orders, OrderSummary{OrderID: order.ID, QuoteID: order.QuoteID, CustomerID: order.CustomerID, Status: order.Status, LineCount: len(order.Lines)})
	}
	return orders
}

func orderDetails(order Order) OrderDetails {
	lines := make([]OrderLineDetails, 0, len(order.Lines))
	for _, line := range order.Lines {
		lines = append(lines, OrderLineDetails{ProductSKU: line.ProductSKU, ProductName: line.ProductName, Quantity: line.Quantity, UnitPrice: line.UnitPrice})
	}
	return OrderDetails{OrderID: order.ID, QuoteID: order.QuoteID, CustomerID: order.CustomerID, Status: order.Status, LineCount: len(lines), ShippedAt: order.ShippedAt, Lines: lines}
}
