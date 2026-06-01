package application

import (
	"testing"
	"time"

	"onion-architecture/internal/domain"
)

func TestGetReturnRequestServiceReturnsDetails(t *testing.T) {
	returns := &stubReturnRequestStore{
		found: domain.ReturnRequest{
			ID:          "return-001",
			OrderID:     "order-001",
			Status:      domain.ReturnRequestStatusRequested,
			Reason:      "damaged on arrival",
			RequestedAt: time.Date(2026, 6, 5, 10, 0, 0, 0, time.UTC),
			RequestedBy: "customer-001",
		},
	}

	service := NewGetReturnRequestService(returns)

	result, err := service.Execute(GetReturnRequestQuery{ReturnRequestID: "return-001"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.ReturnRequestID != "return-001" {
		t.Fatalf("expected return id return-001, got %s", result.ReturnRequestID)
	}

	if result.RequestedBy != "customer-001" {
		t.Fatalf("expected requested by customer-001, got %s", result.RequestedBy)
	}
}

func TestListReturnRequestsServiceFiltersByStatus(t *testing.T) {
	returns := &stubReturnRequestStore{
		list: []domain.ReturnRequest{
			{
				ID:          "return-001",
				OrderID:     "order-001",
				Status:      domain.ReturnRequestStatusRequested,
				RequestedBy: "customer-001",
			},
			{
				ID:          "return-002",
				OrderID:     "order-002",
				Status:      domain.ReturnRequestStatusRefunded,
				RequestedBy: "customer-002",
			},
		},
	}

	service := NewListReturnRequestsService(returns)

	result, err := service.Execute(ListReturnRequestsQuery{Status: domain.ReturnRequestStatusRequested})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}

	if result[0].ReturnRequestID != "return-001" {
		t.Fatalf("expected return-001, got %s", result[0].ReturnRequestID)
	}
}
