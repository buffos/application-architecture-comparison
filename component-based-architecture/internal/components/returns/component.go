package returns

import (
	"component-based-architecture/internal/components/inventory"
	"component-based-architecture/internal/components/orders"
	"component-based-architecture/internal/components/payments"
	"fmt"
)

const (
	ReturnRequestStatusRequested = "Requested"
	ReturnRequestStatusRefunded  = "Refunded"
	ReturnRequestStatusRejected  = "Rejected"
)

type ReturnRequest struct {
	ID, OrderID, CustomerID, Reason, Status string
	LineCount                               int
	amount                                  int
	restock                                 []inventory.RestockItem
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
type ReviewReturnCommand struct{ ReturnRequestID string }

func (c *Component) RequestReturn(command RequestReturnCommand) (RequestReturnResult, error) {
	order, err := c.orders.GetReturnableOrder(command.OrderID)
	if err != nil {
		return RequestReturnResult{}, err
	}
	c.nextID++
	amount := 0
	restock := make([]inventory.RestockItem, 0, len(order.Lines))
	for _, line := range order.Lines {
		amount += line.Quantity * line.UnitPrice
		restock = append(restock, inventory.RestockItem{ProductSKU: line.ProductSKU, Quantity: line.Quantity})
	}
	request := ReturnRequest{ID: fmt.Sprintf("return-%03d", c.nextID), OrderID: order.OrderID, CustomerID: order.CustomerID, Reason: command.Reason, Status: ReturnRequestStatusRequested, LineCount: len(order.Lines), amount: amount, restock: restock}
	c.requests[request.ID] = request
	return RequestReturnResult{ReturnRequestID: request.ID, OrderID: request.OrderID, CustomerID: request.CustomerID, Status: request.Status, LineCount: request.LineCount}, nil
}
func (c *Component) AcceptReturn(command ReviewReturnCommand) error {
	r, ok := c.requests[command.ReturnRequestID]
	if !ok || r.Status != ReturnRequestStatusRequested {
		return fmt.Errorf("return request is not reviewable")
	}
	if err := c.payments.Refund(payments.RefundRequest{OrderID: r.OrderID, CustomerID: r.CustomerID, Amount: r.amount, Reason: r.Reason}); err != nil {
		return err
	}
	if err := c.inventory.Restock(r.restock); err != nil {
		return err
	}
	r.Status = ReturnRequestStatusRefunded
	c.requests[r.ID] = r
	return nil
}
func (c *Component) RejectReturn(command ReviewReturnCommand) error {
	r, ok := c.requests[command.ReturnRequestID]
	if !ok || r.Status != ReturnRequestStatusRequested {
		return fmt.Errorf("return request is not reviewable")
	}
	r.Status = ReturnRequestStatusRejected
	c.requests[r.ID] = r
	return nil
}
