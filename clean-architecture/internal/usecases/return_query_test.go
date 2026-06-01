package usecases

import (
	"testing"

	"clean-architecture/internal/entities"
)

type stubReturnRequestReader struct {
	request entities.ReturnRequest
	err     error
}

func (g stubReturnRequestReader) FindByID(id string) (entities.ReturnRequest, error) {
	if g.err != nil {
		return entities.ReturnRequest{}, g.err
	}

	return g.request, nil
}

type stubReturnRequestLister struct {
	requests []entities.ReturnRequest
	err      error
	status   string
}

func (g *stubReturnRequestLister) ListByStatus(status string) ([]entities.ReturnRequest, error) {
	g.status = status
	if g.err != nil {
		return nil, g.err
	}

	return g.requests, nil
}

type stubGetReturnRequestOutput struct {
	output GetReturnRequestOutput
}

func (o *stubGetReturnRequestOutput) Present(output GetReturnRequestOutput) error {
	o.output = output
	return nil
}

type stubListReturnRequestsOutput struct {
	output ListReturnRequestsOutput
}

func (o *stubListReturnRequestsOutput) Present(output ListReturnRequestsOutput) error {
	o.output = output
	return nil
}

func TestGetReturnRequestInteractorLoadsReturnRequest(t *testing.T) {
	output := &stubGetReturnRequestOutput{}
	interactor := NewGetReturnRequestInteractor(stubReturnRequestReader{
		request: entities.ReturnRequest{
			ID:          "return-001",
			OrderID:     "order-001",
			Reason:      "damaged item",
			Status:      entities.ReturnRequestStatusRequested,
			RequestedBy: "customer-001",
			ReviewedBy:  "reviewer-001",
			ProcessedBy: "finance-001",
		},
	}, output)

	err := interactor.Execute(GetReturnRequestInput{ReturnRequestID: "return-001"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output.output.ReturnRequestID != "return-001" {
		t.Fatalf("expected return id return-001, got %s", output.output.ReturnRequestID)
	}

	if output.output.Status != entities.ReturnRequestStatusRequested {
		t.Fatalf("expected status %s, got %s", entities.ReturnRequestStatusRequested, output.output.Status)
	}
}

func TestListReturnRequestsInteractorFiltersByStatus(t *testing.T) {
	returns := &stubReturnRequestLister{
		requests: []entities.ReturnRequest{
			{
				ID:      "return-001",
				OrderID: "order-001",
				Reason:  "damaged item",
				Status:  entities.ReturnRequestStatusRequested,
			},
			{
				ID:      "return-002",
				OrderID: "order-002",
				Reason:  "changed mind",
				Status:  entities.ReturnRequestStatusRequested,
			},
		},
	}
	output := &stubListReturnRequestsOutput{}
	interactor := NewListReturnRequestsInteractor(returns, output)

	err := interactor.Execute(ListReturnRequestsInput{Status: entities.ReturnRequestStatusRequested})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if returns.status != entities.ReturnRequestStatusRequested {
		t.Fatalf("expected status filter %s, got %s", entities.ReturnRequestStatusRequested, returns.status)
	}

	if output.output.Count != 2 {
		t.Fatalf("expected 2 returns, got %d", output.output.Count)
	}

	if output.output.Requests[0].ReturnRequestID != "return-001" {
		t.Fatalf("expected first return id return-001, got %s", output.output.Requests[0].ReturnRequestID)
	}
}
