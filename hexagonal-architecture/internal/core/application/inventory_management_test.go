package application

import (
	"testing"

	"hexagonal-architecture/internal/adapters/repository/memory"
	"hexagonal-architecture/internal/core/domain"
)

func TestInventoryManagementUseCases(t *testing.T) {
	productRepo := memory.NewProductRepository()
	inventory := memory.NewInventoryReservationAdapter(map[string]int{
		"CHAIR-001": 2,
	})

	_ = productRepo.Save(domain.Product{SKU: "CHAIR-001", Name: "Office Chair", Category: "Standard", BasePrice: 10000, Available: true, ReturnWindowDays: 30})
	_ = productRepo.Save(domain.Product{SKU: "DESK-001", Name: "Executive Desk", Category: "CustomBuild", BasePrice: 50000, Available: true, ReturnWindowDays: 30})

	receiveStock := NewReceiveStockUseCase(productRepo, inventory)
	adjustThreshold := NewAdjustReorderThresholdUseCase(inventory)
	getStockRecord := NewGetStockRecordUseCase(inventory)

	received, err := receiveStock.Execute("CHAIR-001", 3)
	if err != nil {
		t.Fatalf("expected receive stock to succeed, got %v", err)
	}

	if received.Available != 5 {
		t.Fatalf("expected available stock 5, got %d", received.Available)
	}

	newRecord, err := receiveStock.Execute("DESK-001", 4)
	if err != nil {
		t.Fatalf("expected receive stock for new record to succeed, got %v", err)
	}

	if newRecord.Available != 4 {
		t.Fatalf("expected available stock 4, got %d", newRecord.Available)
	}

	adjusted, err := adjustThreshold.Execute("CHAIR-001", 3)
	if err != nil {
		t.Fatalf("expected threshold adjustment to succeed, got %v", err)
	}

	if adjusted.ReorderThreshold != 3 {
		t.Fatalf("expected reorder threshold 3, got %d", adjusted.ReorderThreshold)
	}

	record, err := getStockRecord.Execute("CHAIR-001")
	if err != nil {
		t.Fatalf("expected get stock record to succeed, got %v", err)
	}

	if record.Available != 5 || record.ReorderThreshold != 3 {
		t.Fatalf("unexpected stock record: %+v", record)
	}
}

func TestInventoryManagementValidation(t *testing.T) {
	productRepo := memory.NewProductRepository()
	inventory := memory.NewInventoryReservationAdapter(map[string]int{
		"CHAIR-001": 2,
	})

	_ = productRepo.Save(domain.Product{SKU: "CHAIR-001", Name: "Office Chair", Category: "Standard", BasePrice: 10000, Available: true, ReturnWindowDays: 30})

	receiveStock := NewReceiveStockUseCase(productRepo, inventory)
	adjustThreshold := NewAdjustReorderThresholdUseCase(inventory)

	_, err := receiveStock.Execute("CHAIR-001", 0)
	if err != domain.ErrStockQuantityInvalid {
		t.Fatalf("expected %v, got %v", domain.ErrStockQuantityInvalid, err)
	}

	_, err = adjustThreshold.Execute("CHAIR-001", -1)
	if err != domain.ErrReorderThresholdInvalid {
		t.Fatalf("expected %v, got %v", domain.ErrReorderThresholdInvalid, err)
	}
}
