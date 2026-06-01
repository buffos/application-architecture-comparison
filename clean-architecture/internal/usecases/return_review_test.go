package usecases

import (
	"testing"

	"clean-architecture/internal/entities"
)

type stubReturnRequestEditor struct {
	request entities.ReturnRequest
	err     error
	saved   entities.ReturnRequest
}

func (g *stubReturnRequestEditor) FindByID(id string) (entities.ReturnRequest, error) {
	if g.err != nil {
		return entities.ReturnRequest{}, g.err
	}

	return g.request, nil
}

func (g *stubReturnRequestEditor) Save(request entities.ReturnRequest) error {
	g.saved = request
	return nil
}

type stubAcceptReturnOutput struct {
	output AcceptReturnOutput
}

func (o *stubAcceptReturnOutput) Present(output AcceptReturnOutput) error {
	o.output = output
	return nil
}

type stubRejectReturnOutput struct {
	output RejectReturnOutput
}

func (o *stubRejectReturnOutput) Present(output RejectReturnOutput) error {
	o.output = output
	return nil
}

func TestRequestReturnInteractorCreatesRequestedReturn(t *testing.T) {
	orders := &stubOrderEditor{
		order: entities.Order{
			ID:            "order-001",
			CustomerID:    "customer-001",
			SourceQuoteID: "quote-001",
			Status:        entities.OrderStatusShipped,
			Lines: []entities.OrderLine{
				{SKU: "CHAIR-001", Quantity: 2},
			},
		},
	}
	returns := &stubReturnRequestWriter{}
	output := &stubRequestReturnOutput{}

	interactor := NewRequestReturnInteractor(orders, returns, output)

	err := interactor.Execute(RequestReturnInput{OrderID: "order-001", Reason: "damaged item"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if returns.saved.Status != entities.ReturnRequestStatusRequested {
		t.Fatalf("expected status %s, got %s", entities.ReturnRequestStatusRequested, returns.saved.Status)
	}
}

func TestAcceptReturnInteractorRefundsAndRestocks(t *testing.T) {
	returns := &stubReturnRequestEditor{
		request: entities.ReturnRequest{
			ID:      "return-001",
			OrderID: "order-001",
			Reason:  "damaged item",
			Status:  entities.ReturnRequestStatusRequested,
		},
	}
	orders := &stubOrderEditor{
		order: entities.Order{
			ID:            "order-001",
			CustomerID:    "customer-001",
			SourceQuoteID: "quote-001",
			Status:        entities.OrderStatusShipped,
			Lines: []entities.OrderLine{
				{SKU: "CHAIR-001", Quantity: 2},
			},
		},
	}
	restock := &stubInventoryRestock{}
	output := &stubAcceptReturnOutput{}

	interactor := NewAcceptReturnInteractor(orders, returns, stubRefundGateway{}, restock, output)

	err := interactor.Execute(AcceptReturnInput{ReturnRequestID: "return-001"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if returns.saved.Status != entities.ReturnRequestStatusRefunded {
		t.Fatalf("expected status %s, got %s", entities.ReturnRequestStatusRefunded, returns.saved.Status)
	}

	if len(restock.items) != 1 {
		t.Fatalf("expected 1 restock item, got %d", len(restock.items))
	}
}

func TestRejectReturnInteractorPreventsRefundAndRestock(t *testing.T) {
	returns := &stubReturnRequestEditor{
		request: entities.ReturnRequest{
			ID:      "return-002",
			OrderID: "order-001",
			Reason:  "changed mind",
			Status:  entities.ReturnRequestStatusRequested,
		},
	}
	output := &stubRejectReturnOutput{}

	interactor := NewRejectReturnInteractor(returns, output)

	err := interactor.Execute(RejectReturnInput{ReturnRequestID: "return-002"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if returns.saved.Status != entities.ReturnRequestStatusRejected {
		t.Fatalf("expected status %s, got %s", entities.ReturnRequestStatusRejected, returns.saved.Status)
	}
}
