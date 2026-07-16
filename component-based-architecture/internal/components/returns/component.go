package returns

import (
	"component-based-architecture/internal/components/clock"
	"component-based-architecture/internal/components/idempotency"
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
	ErrRequestedByRequired    = errors.New("return requester is required")
	ErrReviewedByRequired     = errors.New("return reviewer is required")
	ErrProcessedByRequired    = errors.New("return processor is required")
	ErrIdempotencyKeyRequired = errors.New("idempotency key is required")
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
	idempotency idempotency.Store
	requests    map[string]ReturnRequest
	nextID      int
}

func NewComponent(orders orders.ReturnableOrderSource, payments payments.Refunder, inventory inventory.Restocker, eligibility returneligibility.Evaluator, clock clock.Reader, idempotency idempotency.Store) *Component {
	return &Component{orders: orders, payments: payments, inventory: inventory, eligibility: eligibility, clock: clock, idempotency: idempotency, requests: map[string]ReturnRequest{}}
}

type RequestReturnCommand struct{ OrderID, Reason, RequestedBy string }
type RequestReturnResult struct {
	ReturnRequestID, OrderID, CustomerID, Status string
	LineCount                                    int
}
type ReviewReturnCommand struct{ ReturnRequestID, ReviewedBy, ProcessedBy, ReviewNote, IdempotencyKey string }
type ReviewReturnResult struct {
	ReturnRequestID, OrderID, CustomerID, Status string
	LineCount                                    int
}

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
func (c *Component) AcceptReturn(command ReviewReturnCommand) (ReviewReturnResult, error) {
	if result, ok, err := c.replayedResult(command.IdempotencyKey); err != nil || ok {
		return result, err
	}
	if command.ReviewedBy == "" {
		return ReviewReturnResult{}, ErrReviewedByRequired
	}
	r, ok := c.requests[command.ReturnRequestID]
	if !ok || r.Status != ReturnRequestStatusRequested {
		return ReviewReturnResult{}, fmt.Errorf("return request is not reviewable")
	}
	if !c.eligibility.Allows(returneligibility.Review{ShippedAt: r.shippedAt, RequestedAt: r.requestedAt, Lines: r.returnWindows}) {
		r.Status = ReturnRequestStatusRejected
		r.ReviewedBy = command.ReviewedBy
		r.ReviewNote = command.ReviewNote
		c.requests[r.ID] = r
		return c.storeReviewResult(command.IdempotencyKey, r), nil
	}
	if command.ProcessedBy == "" {
		return ReviewReturnResult{}, ErrProcessedByRequired
	}
	if err := c.payments.Refund(payments.RefundRequest{OrderID: r.OrderID, CustomerID: r.CustomerID, Amount: r.amount, Reason: r.Reason}); err != nil {
		return ReviewReturnResult{}, err
	}
	if err := c.inventory.Restock(r.restock); err != nil {
		return ReviewReturnResult{}, err
	}
	r.Status = ReturnRequestStatusRefunded
	r.ReviewedBy = command.ReviewedBy
	r.ProcessedBy = command.ProcessedBy
	r.ReviewNote = command.ReviewNote
	c.requests[r.ID] = r
	return c.storeReviewResult(command.IdempotencyKey, r), nil
}

func returnWindows(lines []orders.ReturnableOrderLine) []returneligibility.ReviewLine {
	windows := make([]returneligibility.ReviewLine, 0, len(lines))
	for _, line := range lines {
		windows = append(windows, returneligibility.ReviewLine{ReturnWindowDays: line.ReturnWindowDays})
	}
	return windows
}
func (c *Component) RejectReturn(command ReviewReturnCommand) (ReviewReturnResult, error) {
	if result, ok, err := c.replayedResult(command.IdempotencyKey); err != nil || ok {
		return result, err
	}
	if command.ReviewedBy == "" {
		return ReviewReturnResult{}, ErrReviewedByRequired
	}
	r, ok := c.requests[command.ReturnRequestID]
	if !ok || r.Status != ReturnRequestStatusRequested {
		return ReviewReturnResult{}, fmt.Errorf("return request is not reviewable")
	}
	r.Status = ReturnRequestStatusRejected
	r.ReviewedBy = command.ReviewedBy
	r.ReviewNote = command.ReviewNote
	c.requests[r.ID] = r
	return c.storeReviewResult(command.IdempotencyKey, r), nil
}

func (c *Component) replayedResult(key string) (ReviewReturnResult, bool, error) {
	if key == "" {
		return ReviewReturnResult{}, false, ErrIdempotencyKeyRequired
	}
	result, ok := c.idempotency.Find(key)
	if !ok {
		return ReviewReturnResult{}, false, nil
	}
	return ReviewReturnResult{ReturnRequestID: result.ReturnRequestID, OrderID: result.OrderID, CustomerID: result.CustomerID, Status: result.Status, LineCount: result.LineCount}, true, nil
}

func (c *Component) storeReviewResult(key string, request ReturnRequest) ReviewReturnResult {
	result := ReviewReturnResult{ReturnRequestID: request.ID, OrderID: request.OrderID, CustomerID: request.CustomerID, Status: request.Status, LineCount: request.LineCount}
	c.idempotency.Save(key, idempotency.Result{ReturnRequestID: result.ReturnRequestID, OrderID: result.OrderID, CustomerID: result.CustomerID, Status: result.Status, LineCount: result.LineCount})
	return result
}
