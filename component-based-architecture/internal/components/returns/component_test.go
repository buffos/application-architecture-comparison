package returns

import (
	"component-based-architecture/internal/components/inventory"
	"component-based-architecture/internal/components/orders"
	"component-based-architecture/internal/components/payments"
	"errors"
	"testing"
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

func (s *inventoryStub) Restock(items []inventory.RestockItem) error { s.items = items; return nil }

func (s *paymentsStub) Refund(request payments.RefundRequest) error { s.request = request; return nil }
func TestRequestReturnRefundsShippedOrder(t *testing.T) {
	p := &paymentsStub{}
	i := &inventoryStub{}
	c := NewComponent(ordersStub{order: orders.ReturnableOrder{OrderID: "order-001", CustomerID: "customer-001", Lines: []orders.ReturnableOrderLine{{ProductSKU: "sku-001", Quantity: 2, UnitPrice: 15000}}}}, p, i)
	r, e := c.RequestReturn(RequestReturnCommand{OrderID: "order-001", Reason: "damaged"})
	if e != nil {
		t.Fatal(e)
	}
	if r.Status != ReturnRequestStatusRefunded || p.request.Amount != 30000 {
		t.Fatalf("unexpected result %+v refund %+v", r, p.request)
	}
	if len(i.items) != 1 || i.items[0].Quantity != 2 {
		t.Fatalf("unexpected restock %+v", i.items)
	}
}
func TestRequestReturnPropagatesNonShippedError(t *testing.T) {
	c := NewComponent(ordersStub{err: orders.ErrOrderNotReturnable}, &paymentsStub{}, &inventoryStub{})
	_, e := c.RequestReturn(RequestReturnCommand{OrderID: "order-001"})
	if !errors.Is(e, orders.ErrOrderNotReturnable) {
		t.Fatalf("got %v", e)
	}
}
