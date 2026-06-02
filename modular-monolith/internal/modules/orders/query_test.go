package orders

import "testing"

type stubQueryRepository struct {
	orders map[string]Order
}

func (r *stubQueryRepository) Save(order Order) error {
	if r.orders == nil {
		r.orders = make(map[string]Order)
	}
	r.orders[order.ID] = order
	return nil
}

func (r *stubQueryRepository) FindByID(id string) (Order, error) {
	order, ok := r.orders[id]
	if !ok {
		return Order{}, ErrOrderNotFound
	}
	return order, nil
}

func (r *stubQueryRepository) ListByStatus(status string) ([]Order, error) {
	list := make([]Order, 0, len(r.orders))
	for _, order := range r.orders {
		if status == "" || order.Status == status {
			list = append(list, order)
		}
	}
	return list, nil
}

func newQueryService(repository Repository) Service {
	return NewService(
		repository,
		stubApprovedQuoteSource{},
		&stubInventoryReserver{},
		&stubPaymentProcessor{},
		&stubShipmentCreator{},
		stubClock{},
	)
}

func TestGetOrderLoadsStoredOrder(t *testing.T) {
	repository := &stubQueryRepository{orders: map[string]Order{}}
	order := Order{
		ID:         "order-001",
		QuoteID:    "quote-001",
		CustomerID: "customer-001",
		Status:     OrderStatusShipped,
		Lines: []OrderLine{
			{ProductSKU: "sku-001", ProductName: "Desk", Quantity: 2, UnitPrice: 15000},
		},
	}
	_ = repository.Save(order)
	service := newQueryService(repository)

	result, err := service.GetOrder(GetOrderQuery{OrderID: order.ID})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.OrderID != order.ID || result.Status != OrderStatusShipped {
		t.Fatalf("expected stored order details to be returned")
	}
}

func TestListOrdersFiltersByStatus(t *testing.T) {
	repository := &stubQueryRepository{orders: map[string]Order{}}
	_ = repository.Save(Order{
		ID:         "order-001",
		QuoteID:    "quote-001",
		CustomerID: "customer-001",
		Status:     OrderStatusPaid,
		Lines:      []OrderLine{{ProductSKU: "sku-001", ProductName: "Desk", Quantity: 2, UnitPrice: 15000}},
	})
	_ = repository.Save(Order{
		ID:         "order-002",
		QuoteID:    "quote-002",
		CustomerID: "customer-002",
		Status:     OrderStatusShipped,
		Lines:      []OrderLine{{ProductSKU: "sku-002", ProductName: "Chair", Quantity: 1, UnitPrice: 5000}},
	})
	service := newQueryService(repository)

	result, err := service.ListOrders(ListOrdersQuery{Status: OrderStatusShipped})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result) != 1 || result[0].Status != OrderStatusShipped {
		t.Fatalf("expected one shipped order, got %+v", result)
	}
}
