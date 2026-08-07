package inventory

import (
	"errors"
	"testing"
)

func TestReturnRestockingServiceReceivesReturnedQuantity(t *testing.T) {
	record, err := NewStockRecord("sku-001", 5, 1)
	if err != nil {
		t.Fatal(err)
	}
	if err := record.Reserve(2); err != nil {
		t.Fatal(err)
	}
	if err := NewReturnRestockingService().RestockAll(
		map[ProductSKU]*StockRecord{"sku-001": &record},
		[]RestockRequest{{SKU: "sku-001", Quantity: 2}},
	); err != nil {
		t.Fatal(err)
	}
	if record.OnHand() != 7 || record.Reserved() != 2 || record.Available() != 5 {
		t.Fatalf("restocked state: on-hand=%d reserved=%d available=%d", record.OnHand(), record.Reserved(), record.Available())
	}
}

func TestReturnRestockingServiceRejectsUnknownAndInvalidRequests(t *testing.T) {
	record, err := NewStockRecord("sku-001", 5, 1)
	if err != nil {
		t.Fatal(err)
	}
	service := NewReturnRestockingService()
	if err := service.RestockAll(map[ProductSKU]*StockRecord{"sku-001": &record}, []RestockRequest{{SKU: "sku-404", Quantity: 1}}); !errors.Is(err, ErrProductSKURequired) {
		t.Fatalf("unknown sku returned %v", err)
	}
	if err := service.RestockAll(map[ProductSKU]*StockRecord{"sku-001": &record}, []RestockRequest{{SKU: "sku-001", Quantity: 0}}); !errors.Is(err, ErrQuantityMustBePositive) {
		t.Fatalf("invalid quantity returned %v", err)
	}
}
