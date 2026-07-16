package returns

import (
	"component-based-architecture/internal/components/idempotency"
	"component-based-architecture/internal/components/inventory"
	"component-based-architecture/internal/components/orders"
	"component-based-architecture/internal/components/payments"
	"component-based-architecture/internal/components/returneligibility"
	"errors"
	"testing"
	"time"
)

type ordersStub struct {
	order orders.ReturnableOrder
	err   error
}

func (s ordersStub) GetReturnableOrder(id string) (orders.ReturnableOrder, error) {
	return s.order, s.err
}

type paymentsStub struct {
	request     payments.RefundRequest
	refundCalls int
}

type inventoryStub struct {
	items        []inventory.RestockItem
	restockCalls int
}

type fixedClock struct{ now time.Time }

func (c fixedClock) Now() time.Time { return c.now }

func (s *inventoryStub) Restock(items []inventory.RestockItem) error {
	s.restockCalls++
	s.items = items
	return nil
}

func (s *paymentsStub) Refund(request payments.RefundRequest) error {
	s.refundCalls++
	s.request = request
	return nil
}
func TestRequestReturnStoresRequestedReturnWithoutSideEffects(t *testing.T) {
	p := &paymentsStub{}
	i := &inventoryStub{}
	c := NewComponent(ordersStub{order: returnableOrder()}, p, i, returneligibility.NewComponent(), fixedClock{now: time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)}, idempotency.NewComponent())
	r, e := c.RequestReturn(RequestReturnCommand{OrderID: "order-001", Reason: "damaged", RequestedBy: "agent-001"})
	if e != nil {
		t.Fatal(e)
	}
	if r.Status != ReturnRequestStatusRequested || p.request.Amount != 0 || len(i.items) != 0 {
		t.Fatalf("unexpected request %+v refund %+v restock %+v", r, p.request, i.items)
	}
	if _, err := c.AcceptReturn(ReviewReturnCommand{ReturnRequestID: r.ReturnRequestID, ReviewedBy: "reviewer-001", ProcessedBy: "processor-001", ReviewNote: "eligible", IdempotencyKey: "accept-001"}); err != nil {
		t.Fatal(err)
	}
	if p.request.Amount != 30000 || len(i.items) != 1 || i.items[0].Quantity != 2 {
		t.Fatalf("unexpected acceptance refund %+v restock %+v", p.request, i.items)
	}
	if request := c.requests[r.ReturnRequestID]; request.RequestedBy != "agent-001" || request.ReviewedBy != "reviewer-001" || request.ProcessedBy != "processor-001" || request.ReviewNote != "eligible" {
		t.Fatalf("unexpected return metadata %+v", request)
	}
}
func TestRequestReturnPropagatesNonShippedError(t *testing.T) {
	c := NewComponent(ordersStub{err: orders.ErrOrderNotReturnable}, &paymentsStub{}, &inventoryStub{}, returneligibility.NewComponent(), fixedClock{}, idempotency.NewComponent())
	_, e := c.RequestReturn(RequestReturnCommand{OrderID: "order-001", RequestedBy: "agent-001"})
	if !errors.Is(e, orders.ErrOrderNotReturnable) {
		t.Fatalf("got %v", e)
	}
}

func TestAcceptReturnRejectsPolicyBlockedRequestWithoutSideEffects(t *testing.T) {
	p := &paymentsStub{}
	i := &inventoryStub{}
	c := NewComponent(ordersStub{order: returnableOrder()}, p, i, returneligibility.NewComponent(), fixedClock{now: time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)}, idempotency.NewComponent())
	r, err := c.RequestReturn(RequestReturnCommand{OrderID: "order-001", Reason: "damaged", RequestedBy: "agent-001"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.AcceptReturn(ReviewReturnCommand{ReturnRequestID: r.ReturnRequestID, ReviewedBy: "reviewer-001", ReviewNote: "outside window", IdempotencyKey: "accept-001"}); err != nil {
		t.Fatal(err)
	}
	if request := c.requests[r.ReturnRequestID]; request.Status != ReturnRequestStatusRejected {
		t.Fatalf("status = %s, want %s", request.Status, ReturnRequestStatusRejected)
	}
	if p.request.Amount != 0 || len(i.items) != 0 {
		t.Fatalf("blocked return had side effects: refund %+v restock %+v", p.request, i.items)
	}
	if request := c.requests[r.ReturnRequestID]; request.ReviewedBy != "reviewer-001" || request.ProcessedBy != "" || request.ReviewNote != "outside window" {
		t.Fatalf("unexpected rejected-return metadata %+v", request)
	}
}

func TestRequestReturnRequiresRequester(t *testing.T) {
	c := NewComponent(ordersStub{order: returnableOrder()}, &paymentsStub{}, &inventoryStub{}, returneligibility.NewComponent(), fixedClock{}, idempotency.NewComponent())
	_, err := c.RequestReturn(RequestReturnCommand{OrderID: "order-001"})
	if !errors.Is(err, ErrRequestedByRequired) {
		t.Fatalf("got %v", err)
	}
}

func TestRejectReturnRecordsReviewerWithoutSideEffects(t *testing.T) {
	p := &paymentsStub{}
	i := &inventoryStub{}
	c := NewComponent(ordersStub{order: returnableOrder()}, p, i, returneligibility.NewComponent(), fixedClock{now: time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)}, idempotency.NewComponent())
	r, err := c.RequestReturn(RequestReturnCommand{OrderID: "order-001", RequestedBy: "agent-001"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.RejectReturn(ReviewReturnCommand{ReturnRequestID: r.ReturnRequestID, ReviewedBy: "reviewer-001", ReviewNote: "item is damaged by customer", IdempotencyKey: "reject-001"}); err != nil {
		t.Fatal(err)
	}
	request := c.requests[r.ReturnRequestID]
	if request.Status != ReturnRequestStatusRejected || request.ReviewedBy != "reviewer-001" || request.ReviewNote != "item is damaged by customer" || request.ProcessedBy != "" {
		t.Fatalf("unexpected rejected return %+v", request)
	}
	if p.request.Amount != 0 || len(i.items) != 0 {
		t.Fatalf("rejection had side effects: refund %+v restock %+v", p.request, i.items)
	}
}

func TestAcceptReturnReplaysStoredResultWithoutSideEffects(t *testing.T) {
	p := &paymentsStub{}
	i := &inventoryStub{}
	c := NewComponent(ordersStub{order: returnableOrder()}, p, i, returneligibility.NewComponent(), fixedClock{now: time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)}, idempotency.NewComponent())
	r, err := c.RequestReturn(RequestReturnCommand{OrderID: "order-001", RequestedBy: "agent-001"})
	if err != nil {
		t.Fatal(err)
	}
	command := ReviewReturnCommand{ReturnRequestID: r.ReturnRequestID, ReviewedBy: "reviewer-001", ProcessedBy: "processor-001", IdempotencyKey: "accept-001"}
	first, err := c.AcceptReturn(command)
	if err != nil {
		t.Fatal(err)
	}
	second, err := c.AcceptReturn(command)
	if err != nil {
		t.Fatal(err)
	}
	if first != second || p.request.Amount != 30000 || len(i.items) != 1 || p.refundCalls != 1 || i.restockCalls != 1 {
		t.Fatalf("unexpected replay: first=%+v second=%+v refund=%+v restock=%+v", first, second, p.request, i.items)
	}
}

func TestRejectReturnReplaysStoredResultWithoutSideEffects(t *testing.T) {
	p := &paymentsStub{}
	i := &inventoryStub{}
	c := NewComponent(ordersStub{order: returnableOrder()}, p, i, returneligibility.NewComponent(), fixedClock{now: time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)}, idempotency.NewComponent())
	r, err := c.RequestReturn(RequestReturnCommand{OrderID: "order-001", RequestedBy: "agent-001"})
	if err != nil {
		t.Fatal(err)
	}
	command := ReviewReturnCommand{ReturnRequestID: r.ReturnRequestID, ReviewedBy: "reviewer-001", IdempotencyKey: "reject-001"}
	first, err := c.RejectReturn(command)
	if err != nil {
		t.Fatal(err)
	}
	second, err := c.RejectReturn(command)
	if err != nil {
		t.Fatal(err)
	}
	if first != second || p.request.Amount != 0 || len(i.items) != 0 {
		t.Fatalf("unexpected replay: first=%+v second=%+v refund=%+v restock=%+v", first, second, p.request, i.items)
	}
}

func TestReviewReturnRequiresIdempotencyKey(t *testing.T) {
	c := NewComponent(ordersStub{order: returnableOrder()}, &paymentsStub{}, &inventoryStub{}, returneligibility.NewComponent(), fixedClock{}, idempotency.NewComponent())
	_, err := c.AcceptReturn(ReviewReturnCommand{})
	if !errors.Is(err, ErrIdempotencyKeyRequired) {
		t.Fatalf("got %v", err)
	}
}

func TestReturnReaderLoadsAndListsReturnRequests(t *testing.T) {
	c := NewComponent(ordersStub{order: returnableOrder()}, &paymentsStub{}, &inventoryStub{}, returneligibility.NewComponent(), fixedClock{now: time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)}, idempotency.NewComponent())
	created, err := c.RequestReturn(RequestReturnCommand{OrderID: "order-001", Reason: "damaged", RequestedBy: "agent-001"})
	if err != nil {
		t.Fatal(err)
	}
	var reader Reader = c
	details, err := reader.GetReturnRequest(GetReturnRequestQuery{ReturnRequestID: created.ReturnRequestID})
	if err != nil {
		t.Fatal(err)
	}
	if details.Reason != "damaged" || details.RequestedBy != "agent-001" || details.Status != ReturnRequestStatusRequested {
		t.Fatalf("unexpected details %+v", details)
	}
	listed := reader.ListReturnRequests(ListReturnRequestsQuery{Status: ReturnRequestStatusRequested})
	if len(listed) != 1 || listed[0].ReturnRequestID != created.ReturnRequestID {
		t.Fatalf("unexpected list %+v", listed)
	}
}

func TestReturnReaderRejectsUnknownRequest(t *testing.T) {
	c := NewComponent(ordersStub{}, &paymentsStub{}, &inventoryStub{}, returneligibility.NewComponent(), fixedClock{}, idempotency.NewComponent())
	_, err := c.GetReturnRequest(GetReturnRequestQuery{ReturnRequestID: "return-999"})
	if !errors.Is(err, ErrReturnRequestNotFound) {
		t.Fatalf("got %v", err)
	}
}

func returnableOrder() orders.ReturnableOrder {
	return orders.ReturnableOrder{
		OrderID: "order-001", CustomerID: "customer-001", ShippedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		Lines: []orders.ReturnableOrderLine{{ProductSKU: "sku-001", Quantity: 2, UnitPrice: 15000, ReturnWindowDays: 30}},
	}
}
