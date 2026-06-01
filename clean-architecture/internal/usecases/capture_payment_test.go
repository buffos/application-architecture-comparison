package usecases

import (
	"testing"

	"clean-architecture/internal/entities"
)

type stubOrderEditor struct {
	order entities.Order
	err   error
	saved entities.Order
}

func (g *stubOrderEditor) FindByID(id string) (entities.Order, error) {
	if g.err != nil {
		return entities.Order{}, g.err
	}

	return g.order, nil
}

func (g *stubOrderEditor) Save(order entities.Order) error {
	g.saved = order
	return nil
}

type stubPaymentGateway struct {
	err error
}

func (g stubPaymentGateway) Capture(order entities.Order) error {
	return g.err
}

type stubCapturePaymentOutput struct {
	output CapturePaymentOutput
}

func (o *stubCapturePaymentOutput) Present(output CapturePaymentOutput) error {
	o.output = output
	return nil
}

func TestCapturePaymentInteractorMarksOrderPaid(t *testing.T) {
	orders := &stubOrderEditor{
		order: entities.Order{
			ID:            "order-001",
			CustomerID:    "customer-001",
			SourceQuoteID: "quote-001",
			Status:        entities.OrderStatusPendingPayment,
			Lines: []entities.OrderLine{
				{SKU: "CHAIR-001", Quantity: 2},
			},
		},
	}
	output := &stubCapturePaymentOutput{}

	interactor := NewCapturePaymentInteractor(orders, stubPaymentGateway{}, output)

	err := interactor.Execute(CapturePaymentInput{OrderID: "order-001"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if orders.saved.Status != entities.OrderStatusPaid {
		t.Fatalf("expected saved status %s, got %s", entities.OrderStatusPaid, orders.saved.Status)
	}

	if output.output.Status != entities.OrderStatusPaid {
		t.Fatalf("expected output status %s, got %s", entities.OrderStatusPaid, output.output.Status)
	}
}

func TestCapturePaymentInteractorRejectsWrongOrderState(t *testing.T) {
	orders := &stubOrderEditor{
		order: entities.Order{
			ID:            "order-002",
			CustomerID:    "customer-001",
			SourceQuoteID: "quote-001",
			Status:        entities.OrderStatusPaid,
		},
	}
	output := &stubCapturePaymentOutput{}

	interactor := NewCapturePaymentInteractor(orders, stubPaymentGateway{}, output)

	err := interactor.Execute(CapturePaymentInput{OrderID: "order-002"})
	if err != entities.ErrQuoteCannotTransition {
		t.Fatalf("expected %v, got %v", entities.ErrQuoteCannotTransition, err)
	}
}
