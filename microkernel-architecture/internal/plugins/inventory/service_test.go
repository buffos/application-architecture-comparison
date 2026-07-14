package inventory

import (
	"testing"

	"microkernel-architecture/internal/kernel"
)

type stubRepository struct {
	stock map[string]StockRecord
}

func (r *stubRepository) FindBySKU(sku string) (StockRecord, error) {
	return r.stock[sku], nil
}

func (r *stubRepository) Save(stock StockRecord) error {
	r.stock[stock.ProductSKU] = stock
	return nil
}

func TestReserve(t *testing.T) {
	repository := &stubRepository{
		stock: map[string]StockRecord{
			"sku-001": {ProductSKU: "sku-001", Available: 10},
		},
	}
	service := NewService(repository)

	err := service.Reserve([]kernel.InventoryReservationItem{
		{ProductSKU: "sku-001", Quantity: 2},
	})
	if err != nil {
		t.Fatalf("expected reservation to succeed, got %v", err)
	}

	if repository.stock["sku-001"].Available != 8 {
		t.Fatalf("expected available stock 8, got %d", repository.stock["sku-001"].Available)
	}
}

func TestReserveRejectsInsufficientStock(t *testing.T) {
	repository := &stubRepository{
		stock: map[string]StockRecord{
			"sku-001": {ProductSKU: "sku-001", Available: 1},
		},
	}
	service := NewService(repository)

	err := service.Reserve([]kernel.InventoryReservationItem{
		{ProductSKU: "sku-001", Quantity: 2},
	})
	if err != ErrInsufficientStock {
		t.Fatalf("expected insufficient stock error, got %v", err)
	}
}

func TestRelease(t *testing.T) {
	repository := &stubRepository{
		stock: map[string]StockRecord{
			"sku-001": {ProductSKU: "sku-001", Available: 8},
		},
	}
	service := NewService(repository)

	err := service.Release([]kernel.InventoryReservationItem{
		{ProductSKU: "sku-001", Quantity: 2},
	})
	if err != nil {
		t.Fatalf("expected release to succeed, got %v", err)
	}

	if repository.stock["sku-001"].Available != 10 {
		t.Fatalf("expected available stock 10, got %d", repository.stock["sku-001"].Available)
	}
}

func TestRestock(t *testing.T) {
	repository := &stubRepository{
		stock: map[string]StockRecord{
			"sku-001": {ProductSKU: "sku-001", Available: 8},
		},
	}
	service := NewService(repository)

	err := service.Restock([]kernel.InventoryReservationItem{
		{ProductSKU: "sku-001", Quantity: 2},
	})
	if err != nil {
		t.Fatalf("expected restock to succeed, got %v", err)
	}

	if repository.stock["sku-001"].Available != 10 {
		t.Fatalf("expected available stock 10, got %d", repository.stock["sku-001"].Available)
	}
}
