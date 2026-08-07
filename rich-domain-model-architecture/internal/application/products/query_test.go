package products

import (
	"errors"
	"testing"

	"rich-domain-model-architecture/internal/domain/catalog"
)

func TestProductQueryProjectsDetailsAndFiltersActiveState(t *testing.T) {
	price, err := catalog.NewPrice(15000, "USD")
	if err != nil {
		t.Fatal(err)
	}
	product, err := catalog.NewProductWithReturnWindow("sku-001", "Desk", catalog.ProductCategoryStandard, price, 30)
	if err != nil {
		t.Fatal(err)
	}
	reader := NewInMemoryReader()
	reader.Save(product)
	details, err := reader.GetProduct("sku-001")
	if err != nil {
		t.Fatal(err)
	}
	if details.BasePriceCents != 15000 || details.ReturnWindowDays != 30 || !details.Active {
		t.Fatalf("details = %+v", details)
	}
	active := true
	rows := reader.ListProducts(&active)
	if len(rows) != 1 || rows[0].SKU != "sku-001" {
		t.Fatalf("rows = %+v", rows)
	}
}

func TestProductQueryReturnsNotFound(t *testing.T) {
	if _, err := NewInMemoryReader().GetProduct("missing"); !errors.Is(err, ErrProductNotFound) {
		t.Fatalf("missing query returned %v", err)
	}
}
