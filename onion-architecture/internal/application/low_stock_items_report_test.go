package application

import (
	"errors"
	"testing"

	"onion-architecture/internal/domain"
)

type stubInventoryStockReader struct {
	list []domain.InventoryStockRecord
	err  error
}

func (r stubInventoryStockReader) ListStock() ([]domain.InventoryStockRecord, error) {
	if r.err != nil {
		return nil, r.err
	}

	return r.list, nil
}

func TestLowStockItemsReportServiceFiltersByThreshold(t *testing.T) {
	service := NewLowStockItemsReportService(stubInventoryStockReader{
		list: []domain.InventoryStockRecord{
			{ProductSKU: "sku-001", Quantity: 12},
			{ProductSKU: "sku-002", Quantity: 5},
			{ProductSKU: "sku-003", Quantity: 2},
		},
	})

	result, err := service.Execute(LowStockItemsReportQuery{Threshold: 5})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 results, got %d", len(result))
	}

	if result[0].ProductSKU != "sku-002" {
		t.Fatalf("expected first low stock item sku-002, got %s", result[0].ProductSKU)
	}

	if result[1].ProductSKU != "sku-003" {
		t.Fatalf("expected second low stock item sku-003, got %s", result[1].ProductSKU)
	}
}

func TestLowStockItemsReportServiceReturnsReaderError(t *testing.T) {
	expectedErr := errors.New("boom")
	service := NewLowStockItemsReportService(stubInventoryStockReader{
		err: expectedErr,
	})

	_, err := service.Execute(LowStockItemsReportQuery{Threshold: 5})
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}
}
