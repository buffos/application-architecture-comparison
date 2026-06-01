package application

import (
	"testing"

	"onion-architecture/internal/domain"
)

type stubOrderFinder struct {
	order domain.Order
	list  []domain.Order
	err   error
}

func (f stubOrderFinder) FindByID(id string) (domain.Order, error) {
	if f.err != nil {
		return domain.Order{}, f.err
	}

	return f.order, nil
}

func (f stubOrderFinder) ListByStatus(status string) ([]domain.Order, error) {
	if f.err != nil {
		return nil, f.err
	}

	result := make([]domain.Order, 0)
	for _, order := range f.list {
		if order.Status == status {
			result = append(result, order)
		}
	}

	return result, nil
}

func TestGetOrderServiceReturnsOrderDetails(t *testing.T) {
	service := NewGetOrderService(stubOrderFinder{
		order: domain.Order{
			ID:         "order-001",
			QuoteID:    "quote-001",
			CustomerID: "customer-001",
			Status:     domain.OrderStatusPaid,
			Lines: []domain.OrderLine{
				{ProductSKU: "sku-002", Quantity: 1},
			},
		},
	})

	result, err := service.Execute(GetOrderQuery{OrderID: "order-001"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.OrderID != "order-001" {
		t.Fatalf("expected order-001, got %s", result.OrderID)
	}

	if result.Status != domain.OrderStatusPaid {
		t.Fatalf("expected status %s, got %s", domain.OrderStatusPaid, result.Status)
	}
}

func TestListOrdersServiceFiltersByStatus(t *testing.T) {
	service := NewListOrdersService(stubOrderFinder{
		list: []domain.Order{
			{ID: "order-001", Status: domain.OrderStatusPaid},
			{ID: "order-002", Status: domain.OrderStatusShipped},
		},
	})

	result, err := service.Execute(ListOrdersQuery{Status: domain.OrderStatusShipped})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}

	if result[0].OrderID != "order-002" {
		t.Fatalf("expected order-002, got %s", result[0].OrderID)
	}
}
