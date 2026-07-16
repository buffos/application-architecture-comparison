package returns

import (
	"component-based-architecture/internal/components/inventory"
	"component-based-architecture/internal/components/orders"
	"component-based-architecture/internal/components/payments"
	"fmt"
)

const ReturnRequestStatusRefunded = "Refunded"

type ReturnRequest struct {
	ID, OrderID, CustomerID, Reason, Status string
	LineCount                               int
}
type Component struct {
	orders    orders.ReturnableOrderSource
	payments  payments.Refunder
	inventory inventory.Restocker
	requests  map[string]ReturnRequest
	nextID    int
}

func NewComponent(orders orders.ReturnableOrderSource, payments payments.Refunder, inventory inventory.Restocker) *Component {
	return &Component{orders: orders, payments: payments, inventory: inventory, requests: map[string]ReturnRequest{}}
}

type RequestReturnCommand struct{ OrderID, Reason string }
type RequestReturnResult struct {
	ReturnRequestID, OrderID, CustomerID, Status string
	LineCount                                    int
}

func (c *Component) RequestReturn(command RequestReturnCommand) (RequestReturnResult, error) {
	order, err := c.orders.GetReturnableOrder(command.OrderID)
	if err != nil {
		return RequestReturnResult{}, err
	}
	amount := 0
	for _, line := range order.Lines {
		amount += line.Quantity * line.UnitPrice
	}
	if err := c.payments.Refund(payments.RefundRequest{OrderID: order.OrderID, CustomerID: order.CustomerID, Amount: amount, Reason: command.Reason}); err != nil {
		return RequestReturnResult{}, err
	}
	restock := make([]inventory.RestockItem, 0, len(order.Lines))
	for _, line := range order.Lines {
		restock = append(restock, inventory.RestockItem{ProductSKU: line.ProductSKU, Quantity: line.Quantity})
	}
	if err := c.inventory.Restock(restock); err != nil {
		return RequestReturnResult{}, err
	}
	c.nextID++
	request := ReturnRequest{ID: fmt.Sprintf("return-%03d", c.nextID), OrderID: order.OrderID, CustomerID: order.CustomerID, Reason: command.Reason, Status: ReturnRequestStatusRefunded, LineCount: len(order.Lines)}
	c.requests[request.ID] = request
	return RequestReturnResult{ReturnRequestID: request.ID, OrderID: request.OrderID, CustomerID: request.CustomerID, Status: request.Status, LineCount: request.LineCount}, nil
}
