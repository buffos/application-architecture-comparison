package inventory

import (
	"errors"
	"testing"
)

func TestInventoryReservationServiceReservesMultipleRecords(t *testing.T) {
	first, err := NewStockRecord("sku-001", 5, 1)
	if err != nil {
		t.Fatal(err)
	}
	second, err := NewStockRecord("sku-002", 3, 1)
	if err != nil {
		t.Fatal(err)
	}
	records := map[ProductSKU]*StockRecord{"sku-001": &first, "sku-002": &second}
	reservations, err := NewInventoryReservationService().ReserveAll(records, []ReservationRequest{{SKU: "sku-001", Quantity: 2}, {SKU: "sku-002", Quantity: 1}})
	if err != nil {
		t.Fatal(err)
	}
	if len(reservations) != 2 || first.Reserved() != 2 || second.Reserved() != 1 {
		t.Fatalf("unexpected reservations %+v stock1=%+v stock2=%+v", reservations, first, second)
	}
}

func TestInventoryReservationServiceRollsBackOnFailure(t *testing.T) {
	first, err := NewStockRecord("sku-001", 5, 1)
	if err != nil {
		t.Fatal(err)
	}
	second, err := NewStockRecord("sku-002", 1, 1)
	if err != nil {
		t.Fatal(err)
	}
	records := map[ProductSKU]*StockRecord{"sku-001": &first, "sku-002": &second}
	_, err = NewInventoryReservationService().ReserveAll(records, []ReservationRequest{{SKU: "sku-001", Quantity: 2}, {SKU: "sku-002", Quantity: 2}})
	if !errors.Is(err, ErrInsufficientAvailable) {
		t.Fatalf("reservation returned %v", err)
	}
	if first.Reserved() != 0 || second.Reserved() != 0 {
		t.Fatalf("rollback failed: stock1=%d stock2=%d", first.Reserved(), second.Reserved())
	}
}
