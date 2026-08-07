package inventory

import (
	"errors"
	"testing"
)

func TestInventoryReservationServiceReservesMultipleAggregates(t *testing.T) {
	first, err := NewStockRecord("sku-001", 10, 2)
	if err != nil {
		t.Fatal(err)
	}
	second, err := NewStockRecord("sku-002", 5, 1)
	if err != nil {
		t.Fatal(err)
	}
	records := map[ProductSKU]*StockRecord{"sku-001": &first, "sku-002": &second}

	reservations, err := NewInventoryReservationService().ReserveAll(records, []ReservationRequest{{SKU: "sku-001", Quantity: 2}, {SKU: "sku-002", Quantity: 1}})
	if err != nil {
		t.Fatal(err)
	}
	if len(reservations) != 2 || first.Reserved() != 2 || second.Reserved() != 1 {
		t.Fatalf("reservations=%+v first=%d second=%d", reservations, first.Reserved(), second.Reserved())
	}
}

func TestInventoryReservationServiceRollsBackEarlierReservations(t *testing.T) {
	first, err := NewStockRecord("sku-001", 10, 2)
	if err != nil {
		t.Fatal(err)
	}
	second, err := NewStockRecord("sku-002", 1, 0)
	if err != nil {
		t.Fatal(err)
	}
	records := map[ProductSKU]*StockRecord{"sku-001": &first, "sku-002": &second}

	_, err = NewInventoryReservationService().ReserveAll(records, []ReservationRequest{{SKU: "sku-001", Quantity: 2}, {SKU: "sku-002", Quantity: 2}})
	if !errors.Is(err, ErrInsufficientAvailable) {
		t.Fatalf("reservation error = %v", err)
	}
	if first.Reserved() != 0 || second.Reserved() != 0 {
		t.Fatalf("rollback left reservations: first=%d second=%d", first.Reserved(), second.Reserved())
	}
}

func TestStockRecordOwnsQuantityInvariants(t *testing.T) {
	if _, err := NewStockRecord("", 1, 0); !errors.Is(err, ErrProductSKURequired) {
		t.Fatalf("missing sku returned %v", err)
	}
	if _, err := NewStockRecord("sku-001", -1, 0); !errors.Is(err, ErrOnHandQuantityNegative) {
		t.Fatalf("negative on-hand returned %v", err)
	}
	record, err := NewStockRecord("sku-001", 2, 1)
	if err != nil {
		t.Fatal(err)
	}
	if err := record.Reserve(3); !errors.Is(err, ErrInsufficientAvailable) {
		t.Fatalf("over-reservation returned %v", err)
	}
	if err := record.Receive(0); !errors.Is(err, ErrQuantityMustBePositive) {
		t.Fatalf("zero receive returned %v", err)
	}
}
