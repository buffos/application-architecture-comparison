package catalog

import (
	"errors"
	"testing"
)

func TestProductAggregateOwnsCatalogLifecycle(t *testing.T) {
	price, err := NewPrice(15000, "usd")
	if err != nil {
		t.Fatal(err)
	}
	product, err := NewProduct("sku-001", "Desk", ProductCategoryStandard, price, 30)
	if err != nil {
		t.Fatal(err)
	}
	if err := product.EnsureSellable(); err != nil {
		t.Fatal(err)
	}
	if err := product.Deactivate(); err != nil {
		t.Fatal(err)
	}
	if err := product.EnsureSellable(); !errors.Is(err, ErrProductInactive) {
		t.Fatalf("inactive product returned %v", err)
	}
	if err := product.Deactivate(); !errors.Is(err, ErrProductAlreadyInactive) {
		t.Fatalf("repeated deactivation returned %v", err)
	}
	if err := product.Activate(); err != nil {
		t.Fatal(err)
	}
	updatedPrice, err := NewPrice(17000, "USD")
	if err != nil {
		t.Fatal(err)
	}
	if err := product.ChangeBasePrice(updatedPrice); err != nil {
		t.Fatal(err)
	}
	if product.BasePrice().Cents() != 17000 {
		t.Fatalf("price = %d, want 17000", product.BasePrice().Cents())
	}
}

func TestProductAggregateRejectsInvalidCatalogData(t *testing.T) {
	price, err := NewPrice(15000, "USD")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := NewProduct("", "Desk", ProductCategoryStandard, price, 30); !errors.Is(err, ErrProductSKURequired) {
		t.Fatalf("empty sku returned %v", err)
	}
	if _, err := NewProduct("sku-001", "", ProductCategoryStandard, price, 30); !errors.Is(err, ErrProductNameRequired) {
		t.Fatalf("empty name returned %v", err)
	}
	if _, err := NewProduct("sku-001", "Desk", ProductCategory("Unknown"), price, 30); !errors.Is(err, ErrProductCategoryInvalid) {
		t.Fatalf("invalid category returned %v", err)
	}
	if _, err := NewProduct("sku-001", "Desk", ProductCategoryStandard, price, 0); !errors.Is(err, ErrReturnWindowInvalid) {
		t.Fatalf("invalid return window returned %v", err)
	}
}
