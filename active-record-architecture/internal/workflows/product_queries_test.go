package workflows

import (
	"testing"

	"active-record-architecture/internal/records"
)

func TestGetProductReturnsProductSnapshot(t *testing.T) {
	db := records.NewDatabase()
	product := records.NewProduct(db, "sku-001", "Desk", "Standard", true, 15000)
	if err := product.Save(); err != nil {
		t.Fatalf("Product.Save() error = %v", err)
	}

	got, err := records.GetProduct(db, product.SKU)
	if err != nil {
		t.Fatalf("GetProduct() error = %v", err)
	}
	got.Name = "Changed locally"
	saved, err := records.FindProduct(db, product.SKU)
	if err != nil {
		t.Fatalf("FindProduct() error = %v", err)
	}
	if saved.Name != "Desk" || saved.Category != "Standard" {
		t.Fatalf("stored product = %#v, want original snapshot", saved)
	}
}

func TestListProductsFiltersAndSorts(t *testing.T) {
	db := records.NewDatabase()
	products := []*records.Product{
		records.NewProduct(db, "sku-002", "Inactive Desk", "Standard", false, 15000),
		records.NewProduct(db, "sku-001", "Active Desk", "Standard", true, 15000),
		records.NewProduct(db, "sku-003", "Clearance Lamp", "Clearance", true, 5000),
	}
	for _, product := range products {
		if err := product.Save(); err != nil {
			t.Fatalf("Product.Save() error = %v", err)
		}
	}

	standard, err := records.ListProducts(db, "Standard", true)
	if err != nil {
		t.Fatalf("ListProducts() filtered error = %v", err)
	}
	if len(standard) != 1 || standard[0].SKU != "sku-001" {
		t.Fatalf("active standard products = %#v, want sku-001", standard)
	}

	allStandard, err := records.ListProducts(db, "Standard", false)
	if err != nil {
		t.Fatalf("ListProducts() unrestricted error = %v", err)
	}
	if len(allStandard) != 2 || allStandard[0].SKU != "sku-001" || allStandard[1].SKU != "sku-002" {
		t.Fatalf("all standard products = %#v, want sorted SKUs", allStandard)
	}
}

func TestGetProductRejectsMissingSKU(t *testing.T) {
	db := records.NewDatabase()
	if _, err := records.GetProduct(db, "sku-404"); err != records.ErrProductNotFound {
		t.Fatalf("error = %v, want %v", err, records.ErrProductNotFound)
	}
}
