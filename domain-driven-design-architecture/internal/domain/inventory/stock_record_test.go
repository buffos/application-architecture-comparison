package inventory

import (
	"errors"
	"testing"
)

func TestStockRecordOwnsReservationAndAvailability(t *testing.T) {
	stock, err := NewStockRecord("sku-001", 10, 3)
	if err != nil {
		t.Fatal(err)
	}
	if stock.Available() != 10 || stock.IsLowStock() {
		t.Fatalf("unexpected initial stock: available=%d low=%t", stock.Available(), stock.IsLowStock())
	}
	if err := stock.Reserve(8); err != nil {
		t.Fatal(err)
	}
	if stock.Reserved() != 8 || stock.Available() != 2 || !stock.IsLowStock() {
		t.Fatalf("unexpected reserved stock: %+v", stock)
	}
	if err := stock.Release(3); err != nil {
		t.Fatal(err)
	}
	if stock.Reserved() != 5 || stock.Available() != 5 {
		t.Fatalf("unexpected released stock: %+v", stock)
	}
	if err := stock.Receive(2); err != nil {
		t.Fatal(err)
	}
	if stock.OnHand() != 12 || stock.Available() != 7 {
		t.Fatalf("unexpected received stock: %+v", stock)
	}
}

func TestStockRecordRejectsInvalidMutations(t *testing.T) {
	if _, err := NewStockRecord("sku-001", -1, 3); !errors.Is(err, ErrOnHandQuantityNegative) {
		t.Fatalf("negative on-hand returned %v", err)
	}
	if _, err := NewStockRecord("sku-001", 10, -1); !errors.Is(err, ErrReorderThresholdNegative) {
		t.Fatalf("negative threshold returned %v", err)
	}
	stock, err := NewStockRecord("sku-001", 2, 1)
	if err != nil {
		t.Fatal(err)
	}
	if err := stock.Reserve(3); !errors.Is(err, ErrInsufficientAvailable) {
		t.Fatalf("over-reservation returned %v", err)
	}
	if err := stock.Release(1); !errors.Is(err, ErrReleaseExceedsReserved) {
		t.Fatalf("over-release returned %v", err)
	}
	if err := stock.Receive(0); !errors.Is(err, ErrQuantityMustBePositive) {
		t.Fatalf("zero receipt returned %v", err)
	}
}
