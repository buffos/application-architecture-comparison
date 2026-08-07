package scripts

import (
	"testing"

	"transaction-script-architecture/internal/data"
)

func TestGetProductReturnsProductSnapshot(t *testing.T) {
	store := data.NewStore()
	store.Products["sku-001"] = data.Product{SKU: "sku-001", Name: "Desk", Category: "Standard", Active: true}

	got, err := GetProduct(store, "sku-001")
	if err != nil {
		t.Fatalf("GetProduct() error = %v", err)
	}
	if got.Name != "Desk" || got.Category != "Standard" {
		t.Fatalf("product = %#v, want desk standard product", got)
	}
}

func TestListProductsFiltersAndSorts(t *testing.T) {
	store := data.NewStore()
	store.Products["sku-002"] = data.Product{SKU: "sku-002", Category: "Standard", Active: false}
	store.Products["sku-001"] = data.Product{SKU: "sku-001", Category: "Standard", Active: true}
	store.Products["sku-003"] = data.Product{SKU: "sku-003", Category: "Clearance", Active: true}

	got, err := ListProducts(store, "Standard", true)
	if err != nil {
		t.Fatalf("ListProducts() error = %v", err)
	}
	if len(got) != 1 || got[0].SKU != "sku-001" {
		t.Fatalf("products = %#v, want active sku-001", got)
	}
}

func TestGetProductRejectsMissingProduct(t *testing.T) {
	store := data.NewStore()
	if _, err := GetProduct(store, "sku-404"); err != ErrProductNotFound {
		t.Fatalf("error = %v, want %v", err, ErrProductNotFound)
	}
}
