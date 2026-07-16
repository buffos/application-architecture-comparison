package inventory

import "testing"

func TestReserveRejectsCombinedDuplicateItemsBeyondAvailableStock(t *testing.T) {
	component := NewComponent()
	component.RegisterStock(StockRecord{ProductSKU: "sku-001", Available: 1})

	err := component.Reserve([]ReservationItem{
		{ProductSKU: "sku-001", Quantity: 1},
		{ProductSKU: "sku-001", Quantity: 1},
	})
	if err != ErrInsufficientStock {
		t.Fatalf("expected %v, got %v", ErrInsufficientStock, err)
	}
	if component.stock["sku-001"] != 1 {
		t.Fatalf("expected stock to remain 1, got %d", component.stock["sku-001"])
	}
}
