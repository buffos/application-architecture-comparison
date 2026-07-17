package products

import "testing"

func TestProductReaderLoadsAndFiltersProducts(t *testing.T) {
	component := NewComponent()
	if err := component.Register(Product{SKU: "sku-001", Name: "Desk", Category: "Standard", Active: true, UnitPrice: 15000, ReturnWindowDays: 30}); err != nil {
		t.Fatal(err)
	}
	if err := component.Register(Product{SKU: "sku-002", Name: "Archived Desk", Category: "Standard", Active: false, UnitPrice: 12000}); err != nil {
		t.Fatal(err)
	}
	var reader Reader = component
	details, err := reader.GetProduct(GetProductQuery{SKU: "sku-001"})
	if err != nil || details.ReturnWindowDays != 30 {
		t.Fatalf("details=%+v err=%v", details, err)
	}
	active := true
	listed := reader.ListProducts(ListProductsQuery{Category: "Standard", Active: &active})
	if len(listed) != 1 || listed[0].SKU != "sku-001" {
		t.Fatalf("unexpected list %+v", listed)
	}
}
