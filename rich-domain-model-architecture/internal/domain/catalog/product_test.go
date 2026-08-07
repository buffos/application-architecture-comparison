package catalog

import (
	"errors"
	"testing"
)

func TestProductOwnsCatalogLifecycle(t *testing.T) {
	price, err := NewPrice(15000, "usd")
	if err != nil {
		t.Fatal(err)
	}
	product, err := NewProduct("sku-001", "Desk", ProductCategoryStandard, price)
	if err != nil {
		t.Fatal(err)
	}

	if err := product.EnsureSellable(); err != nil {
		t.Fatal(err)
	}
	updatedPrice, err := NewPrice(17500, "USD")
	if err != nil {
		t.Fatal(err)
	}
	if err := product.ChangePrice(updatedPrice); err != nil {
		t.Fatal(err)
	}
	if product.BasePrice().Cents() != 17500 {
		t.Fatalf("price = %d, want 17500", product.BasePrice().Cents())
	}

	if err := product.Discontinue(); err != nil {
		t.Fatal(err)
	}
	if product.Active() {
		t.Fatal("discontinued product is still active")
	}
	if err := product.EnsureSellable(); !errors.Is(err, ErrProductInactive) {
		t.Fatalf("sellability check returned %v", err)
	}
	if err := product.Discontinue(); !errors.Is(err, ErrProductAlreadyDiscontinued) {
		t.Fatalf("repeated discontinuation returned %v", err)
	}
}

func TestProductRejectsInvalidState(t *testing.T) {
	price, err := NewPrice(1000, "USD")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := NewProduct("", "Desk", ProductCategoryStandard, price); !errors.Is(err, ErrProductSKURequired) {
		t.Fatalf("missing sku returned %v", err)
	}
	if _, err := NewProduct("sku-001", " ", ProductCategoryStandard, price); !errors.Is(err, ErrProductNameRequired) {
		t.Fatalf("missing name returned %v", err)
	}
	if _, err := NewProduct("sku-001", "Desk", ProductCategory("Unknown"), price); !errors.Is(err, ErrProductCategoryInvalid) {
		t.Fatalf("invalid category returned %v", err)
	}
	if _, err := NewPrice(-1, "USD"); !errors.Is(err, ErrPriceAmountNegative) {
		t.Fatalf("negative price returned %v", err)
	}
	if _, err := NewProduct("sku-001", "Desk", ProductCategoryStandard, Price{}); !errors.Is(err, ErrPriceCurrencyRequired) {
		t.Fatalf("empty price returned %v", err)
	}
}
