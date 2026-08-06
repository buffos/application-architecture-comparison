package inventory

import (
	"errors"
	"testing"
)

func TestReturnRestockingServiceReceivesReturnedUnits(t *testing.T) {
	stock, err := NewStockRecord("sku-001", 10, 3)
	if err != nil {
		t.Fatal(err)
	}
	if err := stock.Reserve(4); err != nil {
		t.Fatal(err)
	}
	records := map[ProductSKU]*StockRecord{"sku-001": &stock}
	if err := NewReturnRestockingService().RestockAll(records, []RestockRequest{{SKU: "sku-001", Quantity: 2}}); err != nil {
		t.Fatal(err)
	}
	if stock.OnHand() != 12 || stock.Reserved() != 4 || stock.Available() != 8 {
		t.Fatalf("unexpected restocked stock: %+v", stock)
	}
}

func TestReturnRestockingServiceRejectsInvalidRequest(t *testing.T) {
	stock, err := NewStockRecord("sku-001", 10, 3)
	if err != nil {
		t.Fatal(err)
	}
	records := map[ProductSKU]*StockRecord{"sku-001": &stock}
	if err := NewReturnRestockingService().RestockAll(records, []RestockRequest{{SKU: "sku-001", Quantity: 0}}); !errors.Is(err, ErrQuantityMustBePositive) {
		t.Fatalf("zero restock returned %v", err)
	}
}
