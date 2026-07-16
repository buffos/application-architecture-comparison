package returns

import (
	"component-based-architecture/internal/components/clock"
	"component-based-architecture/internal/components/inventory"
	"component-based-architecture/internal/components/orders"
	"component-based-architecture/internal/components/payments"
	"component-based-architecture/internal/components/returneligibility"
	"errors"
	"fmt"
	"time"
)

const (
	ReturnRequestStatusRequested = "Requested"
	ReturnRequestStatusRefunded  = "Refunded"
	ReturnRequestStatusRejected  = "Rejected"
)

var (
	ErrRequestedByRequired = errors.New("return requester is required")
	ErrReviewedByRequired  = errors.New("return reviewer is required")
	ErrProcessedByRequired = errors.New("return processor is required")
)

type ReturnRequest struct {
	ID, OrderID, CustomerID, Reason, Status string
	LineCount                               int
	amount                                  int
	restock                                 []inventory.RestockItem
	shippedAt, requestedAt                  time.Time
	returnWindows                           []returneligibility.ReviewLine
	RequestedBy, ReviewedBy, ProcessedBy    string
	ReviewNote                              string
}
type Component struct {
	orders      orders.ReturnableOrderSource
	payments    payments.Refunder
	inventory   inventory.Restocker
	eligibility returneligibility.Evaluator
	clock       clock.Reader
	requests    map[string]ReturnRequest
	nextID      int
}

func NewComponent(orders orders.ReturnableOrderSource, payments payments.Refunder, inventory inventory.Restocker, eligibility returneligibility.Evaluator, clock clock.Reader) *Component {
	return &Component{orders: orders, payments: payments, inventory: inventory, eligibility: eligibility, clock: clock, requests: map[string]ReturnRequest{}}
}

type RequestReturnCommand struct{ OrderID, Reason, RequestedBy string }
type RequestReturnResult struct {
	ReturnRequestID, OrderID, CustomerID, Status string
	LineCount                                    int
}
type ReviewReturnCommand struct{ ReturnRequestID, ReviewedBy, ProcessedBy, ReviewNote string }

func (c *Component) RequestReturn(command RequestReturnCommand) (RequestReturnResult, error) {
	if command.RequestedBy == "" {
		return RequestReturnResult{}, ErrRequestedByRequired
	}
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
	request := ReturnRequest{ID: fmt.Sprintf("return-%03d", c.nextID), OrderID: order.OrderID, CustomerID: order.CustomerID, Reason: command.Reason, Status: ReturnRequestStatusRequested, LineCount: len(order.Lines), amount: amount, restock: restock, shippedAt: order.ShippedAt, requestedAt: c.clock.Now(), returnWindows: returnWindows(order.Lines), RequestedBy: command.RequestedBy}
	c.requests[request.ID] = request
	return RequestReturnResult{ReturnRequestID: request.ID, OrderID: request.OrderID, CustomerID: request.CustomerID, Status: request.Status, LineCount: request.LineCount}, nil
}
func (c *Component) AcceptReturn(command ReviewReturnCommand) error {
	if command.ReviewedBy == "" {
		return ErrReviewedByRequired
	}
	r, ok := c.requests[command.ReturnRequestID]
	if !ok || r.Status != ReturnRequestStatusRequested {
		return fmt.Errorf("return request is not reviewable")
	}
	if !c.eligibility.Allows(returneligibility.Review{ShippedAt: r.shippedAt, RequestedAt: r.requestedAt, Lines: r.returnWindows}) {
		r.Status = ReturnRequestStatusRejected
		r.ReviewedBy = command.ReviewedBy
		r.ReviewNote = command.ReviewNote
		c.requests[r.ID] = r
		return nil
	}
	if command.ProcessedBy == "" {
		return ErrProcessedByRequired
	}
	if err := c.payments.Refund(payments.RefundRequest{OrderID: r.OrderID, CustomerID: r.CustomerID, Amount: r.amount, Reason: r.Reason}); err != nil {
		return err
	}
	if err := c.inventory.Restock(r.restock); err != nil {
		return err
	}
	r.Status = ReturnRequestStatusRefunded
	r.ReviewedBy = command.ReviewedBy
	r.ProcessedBy = command.ProcessedBy
	r.ReviewNote = command.ReviewNote
	c.requests[r.ID] = r
	return nil
}

func returnWindows(lines []orders.ReturnableOrderLine) []returneligibility.ReviewLine {
	windows := make([]returneligibility.ReviewLine, 0, len(lines))
	for _, line := range lines {
		windows = append(windows, returneligibility.ReviewLine{ReturnWindowDays: line.ReturnWindowDays})
	}
	return windows
}
func (c *Component) RejectReturn(command ReviewReturnCommand) error {
	if command.ReviewedBy == "" {
		return ErrReviewedByRequired
	}
	r, ok := c.requests[command.ReturnRequestID]
	if !ok || r.Status != ReturnRequestStatusRequested {
		return fmt.Errorf("return request is not reviewable")
	}
	r.Status = ReturnRequestStatusRejected
	r.ReviewedBy = command.ReviewedBy
	r.ReviewNote = command.ReviewNote
	c.requests[r.ID] = r
	return nil
}
