package returns

import (
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

type paymentsStub struct{ request payments.RefundRequest }

type inventoryStub struct{ items []inventory.RestockItem }

type fixedClock struct{ now time.Time }

func (c fixedClock) Now() time.Time { return c.now }

func (s *inventoryStub) Restock(items []inventory.RestockItem) error { s.items = items; return nil }

func (s *paymentsStub) Refund(request payments.RefundRequest) error { s.request = request; return nil }
func TestRequestReturnStoresRequestedReturnWithoutSideEffects(t *testing.T) {
	p := &paymentsStub{}
	i := &inventoryStub{}
	c := NewComponent(ordersStub{order: returnableOrder()}, p, i, returneligibility.NewComponent(), fixedClock{now: time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)})
	r, e := c.RequestReturn(RequestReturnCommand{OrderID: "order-001", Reason: "damaged"})
	if e != nil {
		t.Fatal(e)
	}
	if r.Status != ReturnRequestStatusRequested || p.request.Amount != 0 || len(i.items) != 0 {
		t.Fatalf("unexpected request %+v refund %+v restock %+v", r, p.request, i.items)
	}
	if err := c.AcceptReturn(ReviewReturnCommand{ReturnRequestID: r.ReturnRequestID}); err != nil {
		t.Fatal(err)
	}
	if p.request.Amount != 30000 || len(i.items) != 1 || i.items[0].Quantity != 2 {
		t.Fatalf("unexpected acceptance refund %+v restock %+v", p.request, i.items)
	}
}
func TestRequestReturnPropagatesNonShippedError(t *testing.T) {
	c := NewComponent(ordersStub{err: orders.ErrOrderNotReturnable}, &paymentsStub{}, &inventoryStub{}, returneligibility.NewComponent(), fixedClock{})
	_, e := c.RequestReturn(RequestReturnCommand{OrderID: "order-001"})
	if !errors.Is(e, orders.ErrOrderNotReturnable) {
		t.Fatalf("got %v", e)
	}
}

func TestAcceptReturnRejectsPolicyBlockedRequestWithoutSideEffects(t *testing.T) {
	p := &paymentsStub{}
	i := &inventoryStub{}
	c := NewComponent(ordersStub{order: returnableOrder()}, p, i, returneligibility.NewComponent(), fixedClock{now: time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)})
	r, err := c.RequestReturn(RequestReturnCommand{OrderID: "order-001", Reason: "damaged"})
	if err != nil {
		t.Fatal(err)
	}
	if err := c.AcceptReturn(ReviewReturnCommand{ReturnRequestID: r.ReturnRequestID}); err != nil {
		t.Fatal(err)
	}
	if request := c.requests[r.ReturnRequestID]; request.Status != ReturnRequestStatusRejected {
		t.Fatalf("status = %s, want %s", request.Status, ReturnRequestStatusRejected)
	}
	if p.request.Amount != 0 || len(i.items) != 0 {
		t.Fatalf("blocked return had side effects: refund %+v restock %+v", p.request, i.items)
	}
}

func returnableOrder() orders.ReturnableOrder {
	return orders.ReturnableOrder{
		OrderID: "order-001", CustomerID: "customer-001", ShippedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		Lines: []orders.ReturnableOrderLine{{ProductSKU: "sku-001", Quantity: 2, UnitPrice: 15000, ReturnWindowDays: 30}},
	}
}
