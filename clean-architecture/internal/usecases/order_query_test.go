package usecases

import (
	"testing"

	"clean-architecture/internal/entities"
)

type stubOrderReader struct {
	order entities.Order
	err   error
}

func (g stubOrderReader) FindByID(id string) (entities.Order, error) {
	if g.err != nil {
		return entities.Order{}, g.err
	}

	return g.order, nil
}

type stubOrderLister struct {
	orders []entities.Order
	err    error
	status string
}

func (g *stubOrderLister) ListByStatus(status string) ([]entities.Order, error) {
	g.status = status
	if g.err != nil {
		return nil, g.err
	}

	return g.orders, nil
}

type stubGetOrderOutput struct {
	output GetOrderOutput
}

func (o *stubGetOrderOutput) Present(output GetOrderOutput) error {
	o.output = output
	return nil
}

type stubListOrdersOutput struct {
	output ListOrdersOutput
}

func (o *stubListOrdersOutput) Present(output ListOrdersOutput) error {
	o.output = output
	return nil
}

func TestGetOrderInteractorLoadsOrder(t *testing.T) {
	output := &stubGetOrderOutput{}
	interactor := NewGetOrderInteractor(stubOrderReader{
		order: entities.Order{
			ID:            "order-001",
			CustomerID:    "customer-001",
			SourceQuoteID: "quote-001",
			Status:        entities.OrderStatusPaid,
			Lines: []entities.OrderLine{
				{SKU: "CHAIR-001", Quantity: 2},
			},
		},
	}, output)

	err := interactor.Execute(GetOrderInput{OrderID: "order-001"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output.output.OrderID != "order-001" {
		t.Fatalf("expected order id order-001, got %s", output.output.OrderID)
	}

	if output.output.Lines != 1 {
		t.Fatalf("expected 1 line, got %d", output.output.Lines)
	}
}

func TestListOrdersInteractorFiltersByStatus(t *testing.T) {
	orders := &stubOrderLister{
		orders: []entities.Order{
			{
				ID:            "order-001",
				CustomerID:    "customer-001",
				SourceQuoteID: "quote-001",
				Status:        entities.OrderStatusPaid,
			},
			{
				ID:            "order-002",
				CustomerID:    "customer-002",
				SourceQuoteID: "quote-002",
				Status:        entities.OrderStatusPaid,
			},
		},
	}
	output := &stubListOrdersOutput{}
	interactor := NewListOrdersInteractor(orders, output)

	err := interactor.Execute(ListOrdersInput{Status: entities.OrderStatusPaid})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if orders.status != entities.OrderStatusPaid {
		t.Fatalf("expected status filter %s, got %s", entities.OrderStatusPaid, orders.status)
	}

	if output.output.Count != 2 {
		t.Fatalf("expected 2 orders, got %d", output.output.Count)
	}

	if output.output.Orders[0].OrderID != "order-001" {
		t.Fatalf("expected first order id order-001, got %s", output.output.Orders[0].OrderID)
	}
}
