package application

import (
	"testing"

	"onion-architecture/internal/domain"
)

func TestReturnRateByCategoryReportServiceComputesCategoryMetrics(t *testing.T) {
	orders := stubOrderFinder{
		list: []domain.Order{
			{
				ID:     "order-001",
				Status: domain.OrderStatusShipped,
				Lines: []domain.OrderLine{
					{ProductCategory: "CustomBuild", Quantity: 2},
					{ProductCategory: "Standard", Quantity: 1},
				},
			},
			{
				ID:     "order-002",
				Status: domain.OrderStatusShipped,
				Lines: []domain.OrderLine{
					{ProductCategory: "Standard", Quantity: 3},
				},
			},
		},
	}

	returns := &stubReturnRequestStore{
		list: []domain.ReturnRequest{
			{
				ID:      "return-001",
				OrderID: "order-001",
				Status:  domain.ReturnRequestStatusRefunded,
			},
		},
	}

	service := NewReturnRateByCategoryReportService(orders, returns)

	report, err := service.Execute()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(report) != 2 {
		t.Fatalf("expected 2 categories, got %d", len(report))
	}

	var customBuild ReturnRateByCategoryRow
	for _, row := range report {
		if row.Category == "CustomBuild" {
			customBuild = row
		}
	}

	if customBuild.ShippedQuantity != 2 {
		t.Fatalf("expected custom build shipped quantity 2, got %d", customBuild.ShippedQuantity)
	}

	if customBuild.ReturnedQuantity != 2 {
		t.Fatalf("expected custom build returned quantity 2, got %d", customBuild.ReturnedQuantity)
	}
}
